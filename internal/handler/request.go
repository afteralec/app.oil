package handler

import (
	"context"
	"database/sql"
	"log"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request"

	// TODO: Do away with anything that calls for inner request packages outside of the request package
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
		c.Append(header.HXRedirect, route.RequestPath(rid))
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
		c.Append(header.HXRedirect, route.RequestPath(rid))
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
			return c.Render(view.BadRequest, view.Bind(c), layout.Standalone)
		}

		ft := c.Params("field")
		if len(ft) == 0 {
			c.Status(fiber.StatusBadRequest)
			return c.Render(view.BadRequest, view.Bind(c), layout.Standalone)
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

		if !request.IsFieldTypeValid(req.Type, ft) {
			c.Status(fiber.StatusBadRequest)
			return c.Render(view.BadRequest, view.Bind(c), layout.Standalone)
		}

		field, err := qtx.GetRequestFieldByType(context.Background(), query.GetRequestFieldByTypeParams{
			RID:  rid,
			Type: ft,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		// TODO: Vestigial Request Change Request handling
		// changes, err := qtx.ListCurrentRequestChangeRequestsForRequest(context.Background(), field.ID)
		// if err != nil {
		// 	c.Status(fiber.StatusInternalServerError)
		// 	return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		// }

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

		// TODO: Make this just return the main view or redirect to the logic
		if req.Status == request.StatusIncomplete {
			return c.Redirect(route.RequestPath(rid))
		}

		b := view.Bind(c)
		b, err = request.NewBindFieldView(i.Templates, b, request.BindFieldViewParams{
			PID:     pid,
			Request: &req,
			Field:   &field,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		return c.Render(view.RequestField, b, layout.Standalone)
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
			return c.Render(view.BadRequest, view.Bind(c), layout.Standalone)
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
			if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
				c.Status(fiber.StatusForbidden)
				return c.Render(view.Forbidden, view.Bind(c), layout.Standalone)
			}
		}

		fields, err := qtx.ListRequestFieldsForRequest(context.Background(), rid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		fieldmap := request.FieldMap(fields)

		// TODO: Get this into a utility that returns a struct with utilities
		// changes, err := qtx.ListRequestChangeRequestsForRequest(context.Background(), query.ListRequestChangeRequestsForRequestParams{
		// 	RID:    rid,
		// 	Old:    false,
		// 	Locked: false,
		// })
		// if err != nil {
		// 	c.Status(fiber.StatusInternalServerError)
		// 	return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		// }
		// if len(changes) > 1 {
		// 	// TODO: This is a fatal error
		// 	c.Status(fiber.StatusInternalServerError)
		// 	return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		// }

		// TODO: Finish new bind pattern
		b := view.Bind(c)

		b, err = request.BindDialogs(b, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if req.Status == request.StatusIncomplete {
			// TODO: Validate that NextIncompleteField returns something here
			nifo, err := request.NextIncompleteField(req.Type, fieldmap)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			b, err := request.NewBindFieldView(i.Templates, b, request.BindFieldViewParams{
				PID:     pid,
				Request: &req,
				Field:   nifo.Field,
				Last:    nifo.Last,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			return c.Render(view.RequestField, b, layout.Standalone)
		}

		if req.Status == request.StatusInReview && req.RPID == pid {
			// TODO: Validate that NextUnreviewedField returns something here
			nufo, err := request.NextUnreviewedField(req.Type, fieldmap)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}

			if nufo.Field == nil {
				b, err = request.BindOverview(i.Templates, b, request.BindOverviewParams{
					PID:      pid,
					Request:  &req,
					FieldMap: fieldmap,
				})
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}

				overviewfields, err := request.FieldsForOverview(request.FieldsForOverviewParams{
					PID:      pid,
					Request:  &req,
					FieldMap: fieldmap,
				})
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}

				// TODO: Reintroduce this removed logic
				// changes, err := qtx.ListCurrentRequestChangeRequestsForRequest(context.Background(), req.ID)
				// if err != nil {
				// 	c.Status(fiber.StatusInternalServerError)
				// 	return nil
				// }
				//
				// changeMap := make(map[string]query.RequestChangeRequest)
				// for _, change := range changes {
				// 	changeMap[change.Field] = change
				// }

				// cr, err := request.ContentReview(qtx, req)
				// if err != nil {
				// 	c.Status(fiber.StatusInternalServerError)
				// 	return nil
				// }

				// processedoverviewfields := []field.ForOverview{}
				// for _, overviewfield := range overviewfields {
				// change, ok := changeMap[overviewfield.Name]
				// if ok {
				// 	overviewfield.HasChangeRequest = true
				// 	overviewfield.ChangeRequest = request.BindChangeRequest(request.BindChangeRequestParams{
				// 		PID:           pid,
				// 		ChangeRequest: &change,
				// 	})
				// }

				// 	status, ok := cr.Status(overviewfield.Name)
				// 	if ok && status == request.FieldStatusApproved {
				// 		overviewfield.IsApproved = true
				// 	}
				//
				// 	processedoverviewfields = append(processedoverviewfields, overviewfield)
				// }

				b["SummaryFields"] = overviewfields

				return c.Render(view.RequestOverview, b, layout.Page)
			}

			// TODO: Get this into a utility that returns a struct with utilities
			// unlockedchanges, err := qtx.ListRequestChangeRequestsForRequestField(context.Background(), query.ListRequestChangeRequestsForRequestFieldParams{
			// 	RID:    req.ID,
			// 	Field:  nufo.Field,
			// 	Old:    false,
			// 	Locked: false,
			// })
			// if err != nil {
			// 	c.Status(fiber.StatusInternalServerError)
			// 	return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			// }
			// if len(unlockedchanges) > 1 {
			// 	// TODO: This is a fatal error
			// 	c.Status(fiber.StatusInternalServerError)
			// 	return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			// }

			b, err = request.NewBindFieldView(i.Templates, b, request.BindFieldViewParams{
				PID:     pid,
				Request: &req,
				Last:    nufo.Last,
			})
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			return c.Render(view.RequestField, b, layout.Standalone)
		}

		b, err = request.BindOverview(i.Templates, b, request.BindOverviewParams{
			PID:      pid,
			Request:  &req,
			FieldMap: fieldmap,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		overviewfields, err := request.FieldsForOverview(request.FieldsForOverviewParams{
			PID:      pid,
			Request:  &req,
			FieldMap: fieldmap,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		b["SummaryFields"] = overviewfields

		return c.Render(view.RequestOverview, b, layout.Page)
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

		ft := c.Params("field")
		if len(ft) == 0 {
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

		if !request.IsFieldTypeValid(req.Type, ft) {
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

		field, err := qtx.GetRequestFieldByType(context.Background(), query.GetRequestFieldByTypeParams{
			RID:  rid,
			Type: ft,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = request.UpdateField(qtx, request.UpdateFieldParams{
			PID:     pid,
			Request: &req,
			Field:   &field,
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

		ft := c.Params("field")
		if len(ft) == 0 {
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

		field, err := qtx.GetRequestFieldByType(context.Background(), query.GetRequestFieldByTypeParams{
			RID:  rid,
			Type: ft,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Do we need this count if there's only ever one open change request?
		count, err := qtx.CountOpenRequestChangeRequestsForRequestField(context.Background(), field.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Get this in a request utility
		var status string
		if count > 0 {
			status = request.FieldStatusReviewed
		} else {
			status = request.FieldStatusApproved
		}

		if err = qtx.UpdateRequestFieldStatus(context.Background(), query.UpdateRequestFieldStatusParams{
			ID:     field.ID,
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

		ft := c.Params("field")
		if len(ft) == 0 {
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

		if !request.IsFieldTypeValid(req.Type, ft) {
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

		field, err := qtx.GetRequestFieldByType(context.Background(), query.GetRequestFieldByTypeParams{
			RID:  rid,
			Type: ft,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = qtx.CreateOpenRequestChangeRequest(context.Background(), query.CreateOpenRequestChangeRequestParams{
			RFID: field.ID,
			PID:  pid,
			Text: text,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Look into returning a Boost or specific components here

		c.Append(header.HXRefresh, "true")
		return nil
	}
}

func DeleteRequestChangeRequest(i *service.Interfaces) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pid, err := util.GetPID(c)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		// TODO: Bind this so the permission check is the same as the permission required to create change requests
		// TODO: Or make this more granular
		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		id, err := util.GetID(c)
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

		change, err := qtx.GetOpenRequestChangeRequest(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if change.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = qtx.DeleteRequestChangeRequest(context.Background(), change.ID); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		field, err := qtx.GetRequestField(context.Background(), change.RFID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if field.Status == request.FieldStatusReviewed {
			if err := qtx.UpdateRequestFieldStatus(context.Background(), query.UpdateRequestFieldStatusParams{
				ID:     field.ID,
				Status: request.FieldStatusApproved,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append(header.HXRefresh, "true")
		return nil
	}
}

func EditRequestChangeRequest(i *service.Interfaces) fiber.Handler {
	type input struct {
		Text string
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

		// TODO: Bind this so the permission check is the same as the permission required to create change requests
		// TODO: Or make this more granular
		perms, err := util.GetPermissions(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if !perms.HasPermission(player.PermissionReviewCharacterApplications.Name) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		id, err := util.GetID(c)
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

		change, err := qtx.GetOpenRequestChangeRequest(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if change.PID != pid {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = qtx.EditOpenRequestChangeRequest(context.Background(), query.EditOpenRequestChangeRequestParams{
			ID:   change.ID,
			Text: text,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Use something other than refresh here, either Boost or oob swaps

		c.Append(header.HXRefresh, "true")
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
			fields, err := qtx.ListRequestFieldsForRequest(context.Background(), app.Request.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summary, err := request.NewSummaryForQueue(request.NewSummaryForQueueParams{
				Query:               qtx,
				Request:             &app.Request,
				FieldMap:            request.FieldMap(fields),
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
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		summaries := []request.SummaryForQueue{}
		for _, app := range apps {
			if app.Request.PID == pid {
				continue
			}

			fields, err := qtx.ListRequestFieldsForRequest(context.Background(), app.Request.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			fieldmap := request.FieldMap(fields)
			summary, err := request.NewSummaryForQueue(request.NewSummaryForQueueParams{
				Query:               qtx,
				Request:             &app.Request,
				FieldMap:            fieldmap,
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
