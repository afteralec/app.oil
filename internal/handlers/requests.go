package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/internal/util"
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

		pid := c.Locals("pid")

		if pid == nil {
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
			PID:  pid.(int64),
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
		pid := c.Locals("pid")

		if pid == nil {
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
			PID:  pid.(int64),
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
		lpid := c.Locals("pid")
		if lpid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
		}
		pid, ok := lpid.(int64)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
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

		if req.PID != pid {
			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			iperms, ok := lperms.(permissions.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
			}
			if !iperms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
				c.Status(fiber.StatusForbidden)
				return nil
			}
		}

		_, ok = request.FieldsByType[req.Type]
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		comments := []queries.ListCommentsForRequestWithAuthorRow{}

		if req.PID != pid {
			comments, err = qtx.ListCommentsForRequestWithAuthor(context.Background(), rid)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		b := c.Locals(constants.BindName).(fiber.Map)
		if req.Type == request.TypeCharacterApplication {
			app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), rid)
			if err != nil {
				// TODO: This means that a Request was created without content - this is an error
				// We should instead insert a blank content row here, but deal with this later
				if err == sql.ErrNoRows {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			b = request.BindDialogs(b, request.BindDialogsParams{
				Request: &req,
			})

			b = request.BindCharacterApplicationFieldPage(b, request.BindCharacterApplicationFieldPageParams{
				Application: &app,
				Request:     &req,
				Field:       field,
			})
		} else {
			// TODO: This means that there's a request in the database with an invalid type
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		view := request.GetView(req.Type, field)

		b = request.BindRequestFieldPage(b, request.BindRequestFieldPageParams{
			PID:      pid,
			Field:    field,
			Request:  &req,
			Comments: comments,
		})

		if req.Status == request.StatusIncomplete {
			return c.Redirect(routes.RequestPath(req.ID))
		}

		return c.Render(view, b, "layout-request-field")
	}
}

func RequestPage(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
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

		// TODO: Reviewer path
		// If it's a non-owner viewing, just show a summary
		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		// TODO: Plan for the requests/:id handler (PLAYER)
		// 1. If the request is Incomplete, drop them into the flow -
		//    * Pull the first incomplete field and render that page
		// 2. If the request is Ready with no unresolved comments, render the summary
		// 3. If the request is Ready with unresolved comments, show the next field with an unresolved comment
		// 4. If the request is Reviewed, show an intro page to the review, then all the changes required in one view
		// 5. Player Accepts the Review > back to #3, show the fields with unresolved comments

		// TODO: Finish new bind pattern
		b := c.Locals(constants.BindName).(fiber.Map)
		b = request.BindStatus(b, &req)
		b = request.BindViewedBy(b, request.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = request.BindDialogs(b, request.BindDialogsParams{
			Request: &req,
		})

		content, err := request.GetContent(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		switch req.Status {
		case request.StatusIncomplete:
			field, last := request.GetNextIncompleteField(req.Type, content)
			view := request.GetView(req.Type, field)

			label, description := request.GetFieldLabelAndDescription(req.Type, field)
			b["Header"] = label
			b["SubHeader"] = description

			// TODO: Move to this
			b["FieldLabel"] = label
			b["FieldDescription"] = description

			b["RequestFormID"] = request.FormID

			if last {
				b["UpdateButtonText"] = "Finish"
			} else {
				b["UpdateButtonText"] = "Next"
			}

			b["RequestPath"] = routes.RequestPath(req.ID)
			b["Field"] = field

			// TODO: Get bind exceptions into their own extractor
			if field == request.FieldGender && req.Type == request.TypeCharacterApplication {
				b["GenderNonBinary"] = character.GenderNonBinary
				b["GenderFemale"] = character.GenderFemale
				b["GenderMale"] = character.GenderMale

				b["GenderIsNonBinary"] = content["Gender"] == character.GenderNonBinary
				b["GenderIsFemale"] = content["Gender"] == character.GenderFemale
				b["GenderIsMale"] = content["Gender"] == character.GenderMale
			}

			return c.Render(view, b, "layout-request-field-standalone")
		case request.StatusReady:
			b["HeaderStatusIcon"] = request.MakeStatusIcon(req.Status, 36)
			b["RequestTitle"] = request.GetSummaryTitle(req.Type, content)
			b["SummaryFields"] = request.GetSummaryFields(request.GetSummaryFieldsParams{
				PID:     pid,
				Request: &req,
				Content: content,
			})

			return c.Render("views/requests/content/summary", b, "layout-request-summary")
		}

		// TODO: This means that this request has an invalid status
		c.Status(fiber.StatusInternalServerError)
		return c.Render("views/500", c.Locals(constants.BindName), "standalone")
	}
}

func UpdateRequestField(i *shared.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := new(request.UpdateInput)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		lpid := c.Locals("pid")
		if lpid == nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("views/login", c.Locals(constants.BindName), "views/layouts/standalone")
		}
		pid, ok := lpid.(int64)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
		}

		prid := c.Params("id")
		if len(prid) == 0 {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		rid, err := strconv.ParseInt(prid, 10, 64)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
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

		// Handle a status update
		if field == request.FieldStatus {
			if !request.IsStatusValid(in.Status) {
				c.Status(fiber.StatusBadRequest)
				return nil
			}

			lperms := c.Locals("perms")
			if lperms == nil {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			perms, ok := lperms.(permissions.PlayerGranted)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			ok = request.IsStatusUpdateOK(&req, perms, pid, in.Status)
			if !ok {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			if err = qtx.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
				ID:     rid,
				Status: in.Status,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if err = tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			c.Append("HX-Refresh", "true")
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
		Value string
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
