package handlers

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/layouts"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

func NewRequest(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Type string `form:"type"`
	}

	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !request.IsTypeValid(in.Type) {
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

		// TODO: Limit new requests by type

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
		if in.Type == request.TypeCharacterApplication {
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
func NewCharacterApplication(i *shared.Interfaces) fiber.Handler {
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
			Type: request.TypeCharacterApplication,
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

func RequestFieldPage(i *shared.Interfaces) fiber.Handler {
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

		if !request.IsTypeValid(req.Type) {
			// TODO: This means that there's a request with an invalid type in the system
			c.Status(fiber.StatusInternalServerError)
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
		}

		// TODO: Reviewer path
		if req.PID != pid {
			perms, err := util.GetPermissions(c)
			if err != nil {
				c.Status(fiber.StatusForbidden)
				return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
			}

			if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			c.Status(fiber.StatusForbidden)
			return c.Render(views.Forbidden, views.Bind(c), layouts.Standalone)
		}

		if req.Status == request.StatusIncomplete {
			return c.Redirect(routes.RequestPath(rid))
		}

		view := request.GetView(req.Type, field)

		b := views.Bind(c)
		b = request.BindStatus(b, &req)
		b = request.BindViewedBy(b, request.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = request.BindDialogs(b, request.BindDialogsParams{
			Request: &req,
		})

		content, err := request.Content(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		label, description := request.GetFieldLabelAndDescription(req.Type, field)
		b["FieldLabel"] = label
		b["FieldDescription"] = description

		b["RequestFormID"] = request.FormID

		b["UpdateButtonText"] = "Update"
		b["BackLink"] = routes.RequestPath(rid)

		b["RequestFormPath"] = routes.RequestFieldPath(rid, field)
		b["Field"] = field

		// TODO: Validate this? i.e., make sure that the content map actually has this in there
		b["FieldValue"] = content[field]

		// TODO: Get bind exceptions into their own extractor
		if field == request.FieldGender && req.Type == request.TypeCharacterApplication {
			b["GenderNonBinary"] = character.GenderNonBinary
			b["GenderFemale"] = character.GenderFemale
			b["GenderMale"] = character.GenderMale

			b["GenderIsNonBinary"] = content["Gender"] == character.GenderNonBinary
			b["GenderIsFemale"] = content["Gender"] == character.GenderFemale
			b["GenderIsMale"] = content["Gender"] == character.GenderMale
		}

		return c.Render(view, b, layouts.RequestFieldStandalone)
	}
}

func RequestPage(i *shared.Interfaces) fiber.Handler {
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

		// TODO: Reviewer path
		// If it's a non-owner viewing, just show a summary
		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			// TODO: 403 view
			return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
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
		b = request.BindStatus(b, &req)
		b = request.BindViewedBy(b, request.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = request.BindDialogs(b, request.BindDialogsParams{
			Request: &req,
		})

		content, err := request.Content(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		switch req.Status {
		case request.StatusIncomplete:
			field, last := request.GetNextIncompleteField(req.Type, content)
			view := request.GetView(req.Type, field)

			label, description := request.GetFieldLabelAndDescription(req.Type, field)
			b["FieldLabel"] = label
			b["FieldDescription"] = description

			b["RequestFormID"] = request.FormID

			if last {
				b["UpdateButtonText"] = "Finish"
			} else {
				b["UpdateButtonText"] = "Next"
			}

			b["RequestFormPath"] = routes.RequestFieldPath(req.ID, field)
			b["Field"] = field
			b["FieldValue"] = ""

			// TODO: Get bind exceptions into their own extractor
			if field == request.FieldGender && req.Type == request.TypeCharacterApplication {
				b["GenderNonBinary"] = character.GenderNonBinary
				b["GenderFemale"] = character.GenderFemale
				b["GenderMale"] = character.GenderMale

				b["GenderIsNonBinary"] = content["Gender"] == character.GenderNonBinary
				b["GenderIsFemale"] = content["Gender"] == character.GenderFemale
				b["GenderIsMale"] = content["Gender"] == character.GenderMale
			}

			return c.Render(view, b, layouts.RequestFieldStandalone)
		case request.StatusReady:
			b["HeaderStatusIcon"] = request.MakeStatusIcon(request.MakeStatusIconParams{
				Status:      req.Status,
				Size:        "36",
				IncludeText: true,
			})
			b["RequestTitle"] = request.SummaryTitle(req.Type, content)
			b["SummaryFields"] = request.GetSummaryFields(request.GetSummaryFieldsParams{
				PID:     pid,
				Request: &req,
				Content: content,
			})

			return c.Render(views.RequestSummaryFields, b, layouts.RequestSummary)
		}

		// TODO: This means that this request has an invalid status
		c.Status(fiber.StatusInternalServerError)
		return c.Render(views.InternalServerError, views.Bind(c), layouts.Standalone)
	}
}

func UpdateRequestField(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := new(request.UpdateInput)
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

		field, err := in.GetField()
		if err != nil {
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

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !request.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = in.UpdateField(pid, qtx, &req, field); err != nil {
			if err == request.ErrInvalidInput {
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

		if req.Status == request.StatusIncomplete {
			// TODO: Boost this using the same handler logic for the request page?
			c.Append("HX-Refresh", "true")
		} else {
			c.Append("HX-Redirect", routes.RequestPath(rid))
		}

		return nil
	}
}

func UpdateRequestFieldNew(i *shared.Interfaces) fiber.Handler {
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

		if !request.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}
		if !request.IsFieldValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = request.UpdateField(qtx, request.UpdateFieldParams{
			PID:     pid,
			Request: &req,
			Field:   field,
			Value:   in.Value,
		}); err != nil {
			if err == request.ErrInvalidInput {
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

		if req.Status == request.StatusIncomplete {
			// TODO: Boost this using the same handler logic for the request page?
			c.Append("HX-Refresh", "true")
		} else {
			c.Append("HX-Redirect", routes.RequestPath(rid))
		}

		return nil
	}
}

func UpdateRequestStatus(i *shared.Interfaces) fiber.Handler {
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
		case request.StatusReady:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = request.StatusSubmitted
		case request.StatusSubmitted:
			if req.PID == pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = request.StatusInReview
		case request.StatusInReview:
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
				status = request.StatusReviewed
			} else {
				status = request.StatusApproved
			}
		case request.StatusReviewed:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = request.StatusReady
		case request.StatusApproved:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			// TODO: Figure out resolving an approved request
		case request.StatusRejected:
			if req.PID != pid {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = request.StatusArchived
		default:
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = request.UpdateStatus(qtx, request.UpdateStatusParams{
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

func DeleteRequest(i *shared.Interfaces) fiber.Handler {
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
			if req.Status != request.StatusSubmitted {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = request.StatusRejected
		} else {
			if req.Status == request.StatusArchived || req.Status == request.StatusCanceled {
				// TODO: Figure out deleting an archived or canceled request
				c.Status(fiber.StatusForbidden)
				return nil
			}

			status = request.StatusCanceled
		}

		if err = request.UpdateStatus(qtx, request.UpdateStatusParams{
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
