package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/requests"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

func NewRequest(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Type string `form:"type"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !requests.IsTypeValid(in.Type) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		// TODO: Limit number of new requests by type

		result, err := qtx.CreateRequest(context.Background(), queries.CreateRequestParams{
			PID:  pid,
			Type: in.Type,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		rid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Rework this so there can't be a missing case
		if in.Type == requests.TypeCharacterApplication {
			if err = qtx.CreateCharacterApplicationContent(context.Background(), rid); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		} else {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", routes.RequestPath(rid))
		return nil
	}
}

// TODO: Combine this functionality with the above so it's consistent
func NewCharacterApplication(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		// TODO: Limit new requests by type

		result, err := qtx.CreateRequest(context.Background(), queries.CreateRequestParams{
			PID:  pid,
			Type: requests.TypeCharacterApplication,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		rid, err := result.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = qtx.CreateCharacterApplicationContent(context.Background(), rid); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", routes.RequestPath(rid))
		return nil
	}
}

func RequestFieldPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return c.Render(views.Login, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		field := c.Params("field")
		if len(field) == 0 {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !requests.IsTypeValid(req.Type) {
			// TODO: This means that there's a request with an invalid type in the system
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if !requests.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		content, err := requests.Content(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if req.PID != pid {
			perms, err := util.GetPermissions(c)
			if err != nil {
				c.Status(fiber.StatusForbidden)
				return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
			}
			if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
			}
		}

		if req.Status == requests.StatusIncomplete {
			return c.Redirect(routes.RequestPath(rid))
		}

		view := requests.GetView(req.Type, field)

		b := views.Bind(c)
		b = requests.BindStatus(b, &req)
		b = requests.BindViewedBy(b, requests.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = requests.BindDialogs(b, requests.BindDialogsParams{
			Request: &req,
		})

		label, description := requests.GetFieldLabelAndDescription(req.Type, field)
		b["FieldLabel"] = label
		b["FieldDescription"] = description

		b["RequestFormID"] = requests.FormID

		b["UpdateButtonText"] = "Update"
		b["BackLink"] = routes.RequestPath(rid)

		b["RequestFormPath"] = routes.RequestFieldPath(rid, field)
		b["Field"] = field

		// TODO: Validate this? i.e., make sure that the content map actually has this in there
		b["FieldValue"] = content[field]

		b = requests.BindGenderRadioGroup(b, requests.BindGenderRadioGroupParams{
			Content: content,
			Name:    "value",
		})

		return c.Render(view, b, layouts.RequestFieldStandalone)
	}
}

func RequestPage(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return c.Render(views.Login, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(views.NotFound, views.Bind(c), layouts.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if req.PID != pid {
			perms, err := util.GetPermissions(c)
			if err != nil {
				c.Status(fiber.StatusForbidden)
				return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
			}
			if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
			}
		}

		// TODO: Plan for the requests/:id handler (PLAYER)
		// 1. If the request is Incomplete, drop them into the flow -
		//    * Pull the first incomplete field and render that page
		// 2. If the request is Ready with no unresolved comments, render the summary
		// 3. If the request is Ready with unresolved comments, show the next field with an unresolved comment
		// 4. If the request is Reviewed, show an intro page to the review, then all the changes required in one view
		// 5. Player Accepts the Review > back to #3, show the fields with unresolved comments

		// TODO: Finish new bind pattern
		b := views.Bind(c)
		b = requests.BindStatus(b, &req)
		b = requests.BindViewedBy(b, requests.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = requests.BindDialogs(b, requests.BindDialogsParams{
			Request: &req,
		})

		content, err := requests.Content(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		if req.Status == requests.StatusIncomplete {
			field, last := requests.GetNextIncompleteField(req.Type, content)
			view := requests.GetView(req.Type, field)

			label, description := requests.GetFieldLabelAndDescription(req.Type, field)
			b["FieldLabel"] = label
			b["FieldDescription"] = description

			b["RequestFormID"] = requests.FormID

			if last {
				b["UpdateButtonText"] = "Finish"
			} else {
				b["UpdateButtonText"] = "Next"
			}

			b["RequestFormPath"] = routes.RequestFieldPath(req.ID, field)
			b["Field"] = field
			b["FieldValue"] = ""

			b = requests.BindGenderRadioGroup(b, requests.BindGenderRadioGroupParams{
				Content: content,
				Name:    "value",
			})

			// TODO: Get bind exceptions into their own extractor
			// if field == requests.FieldGender && req.Type == requests.TypeCharacterApplication {
			// 	b["GenderNonBinary"] = character.GenderNonBinary
			// 	b["GenderFemale"] = character.GenderFemale
			// 	b["GenderMale"] = character.GenderMale
			//
			// 	b["GenderIsNonBinary"] = content["Gender"] == character.GenderNonBinary
			// 	b["GenderIsFemale"] = content["Gender"] == character.GenderFemale
			// 	b["GenderIsMale"] = content["Gender"] == character.GenderMale
			// }

			return c.Render(view, b, layouts.RequestFieldStandalone)
		}

		b["PageHeader"] = fiber.Map{
			"Title": requests.SummaryTitle(req.Type, content),
		}
		// TODO: Look at re-implementing this in the view?
		// b["HeaderStatusIcon"] = request.MakeStatusIcon(request.MakeStatusIconParams{
		// 	Status:      req.Status,
		// 	Size:        "36",
		// 	IncludeText: true,
		// })
		b["SummaryFields"] = requests.GetSummaryFields(requests.GetSummaryFieldsParams{
			PID:     pid,
			Request: &req,
			Content: content,
		})

		return c.Render(views.RequestSummaryFields, b, layouts.Page)
	}
}

func UpdateRequestField(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Value string `form:"value"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return nil
			}

			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rid, err := util.GetID(c)
		if err != nil {
			if err == util.ErrNoID {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		field := c.Params("field")
		if len(field) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !requests.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !requests.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !requests.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = requests.UpdateField(qtx, requests.UpdateFieldParams{
			PID:     pid,
			Request: &req,
			Field:   field,
			Value:   in.Value,
		}); err != nil {
			if err == requests.ErrInvalidInput {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if req.Status == requests.StatusIncomplete {
			// TODO: Boost this using the same handler logic for the request page?
			c.Append("HX-Refresh", "true")
		} else {
			c.Append("HX-Redirect", routes.RequestPath(rid))
		}

		return nil
	}
}

func UpdateRequestStatus(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return nil
			}

			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rid, err := util.GetID(c)
		if err != nil {
			if err == util.ErrNoID {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		var status string
		switch req.Status {
		case requests.StatusReady:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = requests.StatusSubmitted
		case requests.StatusSubmitted:
			if req.PID == pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = requests.StatusInReview
		case requests.StatusInReview:
			if req.PID == pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			count, err := qtx.CountUnresolvedComments(context.Background(), rid)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if count > 0 {
				status = requests.StatusReviewed
			} else {
				status = requests.StatusApproved
			}
		case requests.StatusReviewed:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = requests.StatusReady
		case requests.StatusApproved:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			// TODO: Figure out resolving an approved request
		case requests.StatusRejected:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = requests.StatusArchived
		default:
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = requests.UpdateStatus(qtx, requests.UpdateStatusParams{
			RID:    rid,
			PID:    pid,
			Status: status,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return nil
	}
}

func DeleteRequest(i *interfaces.Shared) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		rid, err := util.GetID(c)
		if err != nil {
			if err == util.ErrNoID {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		var status string

		if req.PID != pid {
			if req.Status != requests.StatusSubmitted {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = requests.StatusRejected
		} else {
			if req.Status == requests.StatusArchived || req.Status == requests.StatusCanceled {
				// TODO: Figure out deleting an archived or canceled request
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = requests.StatusCanceled
		}

		if err = requests.UpdateStatus(qtx, requests.UpdateStatusParams{
			RID:    rid,
			PID:    pid,
			Status: status,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		return nil
	}
}
