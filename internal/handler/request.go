package handler

import (
	"context"
	"database/sql"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func NewRequest(i *service.Interfaces) fiber.Handler {
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

		// TODO: Limit number of new requests by type

		rid, err := request.New(qtx, request.NewParams{
			Type: in.Type,
			PID:  pid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", route.RequestPath(rid))
		return nil
	}
}

func NewCharacterApplication(i *service.Interfaces) fiber.Handler {
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

		rid, err := request.New(qtx, request.NewParams{
			Type: request.TypeCharacterApplication,
			PID:  pid,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusCreated)
		c.Append("HX-Redirect", route.RequestPath(rid))
		return nil
	}
}

func RequestFieldPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return c.Render(view.Login, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		field := c.Params("field")
		if len(field) == 0 {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
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
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if !request.IsFieldNameValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		content, err := request.Content(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if req.PID != pid {
			perms, err := util.GetPermissions(c)
			if err != nil {
				c.Status(fiber.StatusForbidden)
				return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
			}
			if !perms.Permissions[player.PermissionReviewCharacterApplications.Name] {
				c.Status(fiber.StatusForbidden)
				return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
			}
		}

		if req.Status == request.StatusIncomplete {
			return c.Redirect(route.RequestPath(rid))
		}

		v := request.View(req.Type, field)

		b := view.Bind(c)
		b = request.BindStatus(b, &req)
		b = request.BindViewedBy(b, request.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = request.BindDialogs(b, request.BindDialogsParams{
			Request: &req,
		})

		label, description := request.GetFieldLabelAndDescription(req.Type, field)
		b["FieldLabel"] = label
		b["FieldDescription"] = description

		b["RequestFormID"] = request.FormID

		b["UpdateButtonText"] = "Update"
		b["BackLink"] = route.RequestPath(rid)

		b["RequestFormPath"] = route.RequestFieldPath(rid, field)
		b["Field"] = field

		// TODO: Use this ok value
		fieldValue, _ := content.Value(field)
		b["FieldValue"] = fieldValue

		// TODO: Let this bind use the actual content API
		b = request.BindGenderRadioGroup(b, request.BindGenderRadioGroupParams{
			Content: content,
			Name:    "value",
		})

		return c.Render(v, b, layout.RequestFieldStandalone)
	}
}

// TODO: Add a back link here based on the request type and viewer
func RequestPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			if err == util.ErrNoPID {
				c.Status(fiber.StatusUnauthorized)
				return c.Render(view.Login, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		rid, err := util.GetID(c)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			// TODO: 400 view
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		req, err := qtx.GetRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return c.Render(view.NotFound, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if req.PID != pid {
			perms, err := util.GetPermissions(c)
			if err != nil {
				c.Status(fiber.StatusForbidden)
				return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
			}
			if !perms.Permissions[player.PermissionReviewCharacterApplications.Name] {
				c.Status(fiber.StatusForbidden)
				return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
			}
		}

		// TODO: Plan for the requests/:id handler (PLAYER)
		// 1. If the request is Incomplete, drop them into the flow -
		//    * Pull the first incomplete field and render that page
		// 2. If the request is Ready with no unresolved comments, render the summary
		// 3. If the request is Ready with unresolved comments, show the next field with an unresolved comment
		// 4. If the request is Reviewed, show an intro page to the review, then all the changes required in one view
		// 5. Player Accepts the Review > back to #3, show the fields with unresolved comments

		content, err := request.Content(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		// TODO: Finish new bind pattern
		b := view.Bind(c)
		b = request.BindStatus(b, &req)
		b = request.BindViewedBy(b, request.BindViewedByParams{
			Request: &req,
			PID:     pid,
		})
		b = request.BindDialogs(b, request.BindDialogsParams{
			Request: &req,
		})

		if req.Status == request.StatusIncomplete {
			field, last := request.NextIncompleteField(req.Type, content)
			view := request.View(req.Type, field)

			// TODO: Validate that NextIncompleteField returns something here

			label, description := request.GetFieldLabelAndDescription(req.Type, field)
			b["FieldLabel"] = label
			b["FieldDescription"] = description

			b["RequestFormID"] = request.FormID

			if last {
				b["UpdateButtonText"] = "Finish"
			} else {
				b["UpdateButtonText"] = "Next"
			}

			b["RequestFormPath"] = route.RequestFieldPath(req.ID, field)
			b["Field"] = field
			b["FieldValue"] = ""

			b = request.BindGenderRadioGroup(b, request.BindGenderRadioGroupParams{
				Content: content,
				Name:    "value",
			})

			return c.Render(view, b, layout.RequestFieldStandalone)
		}

		if req.Status == request.StatusInReview && req.RPID == pid {
			// TODO: Here, the reviewer is viewing a request they're currently
			cr, err := request.ContentReview(qtx, &req)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			// TODO: Validate that NextUnreviewedField returns something here
			field, last := request.NextUnreviewedField(req.Type, cr)
			v := request.View(req.Type, field)
			value, ok := content.Value(field)
			if !ok {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			openChange := true
			change, err := qtx.GetCurrentRequestChangeRequestForRequestField(context.Background(), query.GetCurrentRequestChangeRequestForRequestFieldParams{
				RID:   rid,
				Field: field,
			})
			if err != nil {
				if err == sql.ErrNoRows {
					openChange = false
				} else {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}
			}

			label, description := request.GetFieldLabelAndDescription(req.Type, field)
			b["FieldLabel"] = label
			b["FieldDescription"] = description

			b["RequestFormID"] = request.FormID

			if last {
				b["UpdateButtonText"] = "Finish"
			} else {
				b["UpdateButtonText"] = "Next"
			}

			b["RequestFormPath"] = route.RequestFieldPath(req.ID, field)
			b["Field"] = field
			b["FieldValue"] = value

			b = request.BindGenderRadioGroup(b, request.BindGenderRadioGroupParams{
				Content: content,
				Name:    "value",
			})

			b["ChangeRequestPath"] = route.RequestChangeRequestPath(req.ID, field)
			b["ActionButtonPath"] = route.RequestFieldStatusPath(rid, field)

			if openChange {
				b["ActionButtonText"] = "Next"
				b["ChangeRequest"] = change
			} else {
				b["ActionButtonText"] = "Approve"
			}

			if err := tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			return c.Render(v, b, layout.RequestFieldStandalone)
		}

		b["PageHeader"] = fiber.Map{
			"Title": request.TitleForSummary(req.Type, content),
		}
		// TODO: Build a utility for this
		b["Status"] = fiber.Map{
			"StatusIcon": request.NewStatusIcon(request.StatusIconParams{Status: req.Status, IconSize: 48, IncludeText: true, TextSize: "text-xl"}),
		}
		summaryFields, err := request.FieldsForSummary(request.FieldsForSummaryParams{
			PID:     pid,
			Request: &req,
			Content: content,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		b["SummaryFields"] = summaryFields

		return c.Render(view.RequestSummaryFields, b, layout.Page)
	}
}

func UpdateRequestField(i *service.Interfaces) fiber.Handler {
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

		if !request.IsFieldNameValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsEditable(&req) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = request.UpdateField(qtx, request.UpdateFieldParams{
			PID:       pid,
			Request:   &req,
			FieldName: field,
			Value:     in.Value,
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
			c.Append("HX-Redirect", route.RequestPath(rid))
		}

		return nil
	}
}

func UpdateRequestStatus(i *service.Interfaces) fiber.Handler {
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

		status, err := request.NextStatus(request.NextStatusParams{
			Query:   qtx,
			Request: &req,
			PID:     pid,
		})
		if err != nil {
			if err == request.ErrNextStatusForbidden {
				c.Status(fiber.StatusForbidden)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
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

		c.Append("HX-Refresh", "true")
		return nil
	}
}

func UpdateRequestFieldStatus(i *service.Interfaces) fiber.Handler {
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

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
			log.Println("No permission")
			c.Status(fiber.StatusForbidden)
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

		if req.Status != request.StatusInReview {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if req.RPID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		currentChangeRequestCount, err := qtx.CountCurrentRequestChangeRequestForRequestField(context.Background(), query.CountCurrentRequestChangeRequestForRequestFieldParams{
			RID:   rid,
			Field: field,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Get this in a request utility
		var status string
		if currentChangeRequestCount > 0 {
			status = request.FieldStatusReviewed
		} else {
			status = request.FieldStatusApproved
		}

		if err = request.UpdateFieldStatus(qtx, request.UpdateFieldStatusParams{
			Request:   &req,
			FieldName: field,
			Status:    status,
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
}

func DeleteRequest(i *service.Interfaces) fiber.Handler {
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

func CreateRequestChangeRequest(i *service.Interfaces) fiber.Handler {
	type input struct {
		Text string `form:"text"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		text := request.SanitizeChangeRequestText(in.Text)
		if !request.IsChangeRequestTextValid(text) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		rid, err := util.GetID(c)
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

		if !request.IsFieldNameValid(req.Type, field) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if req.PID == pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if req.Status != request.StatusInReview {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if req.RPID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = qtx.CreateRequestChangeRequest(context.Background(), query.CreateRequestChangeRequestParams{
			RID:   rid,
			PID:   pid,
			Text:  text,
			Field: field,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		_, err = qtx.GetCurrentRequestChangeRequestForRequestField(context.Background(), query.GetCurrentRequestChangeRequestForRequestFieldParams{
			RID:   rid,
			Field: field,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: HTML return
		return nil
	}
}

// TODO: Move this to the Actor file?
func CharactersPage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			// TODO: Pivot this on ErrNoPID
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			// TODO: Figure out what this should redirect to
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		// TODO: Make this a ListRequestsForPlayerByType query instead
		apps, err := qtx.ListCharacterApplicationsForPlayer(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c))
		}

		// TODO: Get this into a standard API on the request package
		summaries := []request.SummaryForQueue{}
		for _, app := range apps {
			content, err := request.Content(qtx, &app.Request)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summary, err := request.NewSummaryForQueue(request.SummaryForQueueParams{
				Query:               qtx,
				Request:             &app.Request,
				Content:             content,
				PID:                 pid,
				ReviewerPermissions: &perms,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summaries = append(summaries, summary)
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := view.Bind(c)
		b["RequestsPath"] = route.Requests
		b["CharacterApplicationSummaries"] = summaries
		b["HasCharacterApplications"] = len(apps) > 0
		return c.Render(view.Characters, b)
	}
}

func CharacterApplicationsQueuePage(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return c.Render(view.Login, view.Bind(c), layout.Standalone)
		}

		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		tx, err := i.Database.Begin()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		defer tx.Rollback()
		qtx := i.Queries.WithTx(tx)

		// TODO: Make this a "List Open Requests By Type" query
		apps, err := qtx.ListOpenCharacterApplications(context.Background())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c))
		}

		summaries := []request.SummaryForQueue{}
		for _, app := range apps {
			content, err := request.Content(qtx, &app.Request)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summary, err := request.NewSummaryForQueue(request.SummaryForQueueParams{
				Query:               qtx,
				Request:             &app.Request,
				Content:             content,
				PID:                 pid,
				ReviewerPermissions: &perms,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summaries = append(summaries, summary)
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		b := view.Bind(c)
		// TODO: Move this length check down into the template
		b["ThereAreCharacterApplications"] = len(summaries) > 0
		b["CharacterApplicationSummaries"] = summaries
		return c.Render(view.CharacterApplicationQueue, b)
	}
}
