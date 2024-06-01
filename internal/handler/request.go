package handler

import (
	"context"
	"database/sql"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/header"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/service"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

func CreateRequest(i *service.Interfaces) fiber.Handler {
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

func CreateCharacterApplication(i *service.Interfaces) fiber.Handler {
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

		rid, err := util.GetID(c, "id")
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

		bfvp := request.BindFieldViewParams{
			PID:     pid,
			Request: &req,
			Field:   &field,
		}
		if request.FieldRequiresSubfields(req.Type, field.Type) {
			subfields, err := qtx.ListRequestSubfieldsForField(context.Background(), field.ID)
			if err != nil {
				if err == sql.ErrNoRows {
				} else {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}
			}
			bfvp.Subfields = subfields
		}

		b := view.Bind(c)
		b, err = request.BindFieldView(i.Templates, b, bfvp)
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

		rid, err := util.GetID(c, "id")
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

		// TODO: Finish new bind pattern
		b := view.Bind(c)
		b, err = request.BindDialogs(b, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		if req.Status == request.StatusIncomplete {
			nifo, err := request.NextIncompleteField(req.Type, fieldmap)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			if nifo.Field == nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			bfvp := request.BindFieldViewParams{
				PID:     pid,
				Request: &req,
				Field:   nifo.Field,
				Last:    nifo.Last,
			}
			if request.FieldRequiresSubfields(req.Type, nifo.Field.Type) {
				subfields, err := qtx.ListRequestSubfieldsForField(context.Background(), nifo.Field.ID)
				if err != nil {
					if err == sql.ErrNoRows {
					} else {
						c.Status(fiber.StatusInternalServerError)
						return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
					}
				}
				bfvp.Subfields = subfields
			}
			b, err := request.BindFieldView(i.Templates, b, bfvp)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			if err := tx.Commit(); err != nil {
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

				// TODO: Extract this to a utility?
				rfids := []int64{}
				for _, field := range fieldmap {
					rfids = append(rfids, field.ID)
				}
				openchanges, err := i.Queries.ListOpenRequestChangeRequestsByFieldID(context.Background(), rfids)
				if err != nil {
					if err == sql.ErrNoRows {
						// TODO: Acceptable, this means that there are no change requests for those fields
					} else {
						c.Status(fiber.StatusInternalServerError)
						return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
					}
				}
				openchangemap := map[int64]query.OpenRequestChangeRequest{}
				for _, change := range openchanges {
					openchangemap[change.RFID] = change
				}
				changes, err := i.Queries.ListRequestChangeRequestsByFieldID(context.Background(), rfids)
				if err != nil {
					if err == sql.ErrNoRows {
						// TODO: Acceptable, this means that there are no change requests for those fields
					} else {
						c.Status(fiber.StatusInternalServerError)
						return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
					}
				}
				changemap := map[int64]query.RequestChangeRequest{}
				for _, change := range changes {
					changemap[change.RFID] = change
				}

				subfields, err := qtx.ListRequestSubfieldsForFields(context.Background(), rfids)
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}

				overviewfields, err := request.FieldsForOverview(i.Templates, request.FieldsForOverviewParams{
					PID:           pid,
					Request:       &req,
					FieldMap:      fieldmap,
					OpenChangeMap: openchangemap,
					ChangeMap:     changemap,
					Subfields:     subfields,
				})
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}

				b["Fields"] = overviewfields

				if err := tx.Commit(); err != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}
				return c.Render(view.RequestOverview, b, layout.Page)
			}

			och, err := i.Queries.GetOpenRequestChangeRequestForRequestField(context.Background(), nufo.Field.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					// TODO: This just means there's no Open Change Request for this field
				} else {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}
			}
			var openchange *query.OpenRequestChangeRequest = nil
			if och.ID != 0 {
				openchange = &och
			}

			// TODO: Rename this query
			ch, err := i.Queries.GetRequestChangeRequestByFieldID(context.Background(), nufo.Field.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					// TODO: This just means there's no Open Change Request for this field
				} else {
					c.Status(fiber.StatusInternalServerError)
					return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
				}
			}
			var change *query.RequestChangeRequest = nil
			if ch.ID != 0 {
				change = &ch
			}

			bfvp := request.BindFieldViewParams{
				PID:        pid,
				Request:    &req,
				Field:      nufo.Field,
				Last:       nufo.Last,
				OpenChange: openchange,
				Change:     change,
			}
			if request.FieldRequiresSubfields(req.Type, nufo.Field.Type) {
				subfields, err := i.Queries.ListRequestSubfieldsForField(context.Background(), nufo.Field.ID)
				if err != nil {
					if err == sql.ErrNoRows {
					} else {
						c.Status(fiber.StatusInternalServerError)
						return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
					}
				}
				bfvp.Subfields = subfields
			}

			b, err = request.BindFieldView(i.Templates, b, bfvp)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}

			if err := tx.Commit(); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			return c.Render(view.RequestField, b, layout.Standalone)
		}

		// TODO: Extract this to a utility?
		rfids := []int64{}
		for _, field := range fieldmap {
			rfids = append(rfids, field.ID)
		}
		changes, err := i.Queries.ListRequestChangeRequestsByFieldID(context.Background(), rfids)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Acceptable, this means that there are no change requests for those fields
			} else {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
		}
		changemap := map[int64]query.RequestChangeRequest{}
		for _, change := range changes {
			changemap[change.RFID] = change
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

		subfields, err := qtx.ListRequestSubfieldsForFields(context.Background(), rfids)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		overviewfields, err := request.FieldsForOverview(i.Templates, request.FieldsForOverviewParams{
			PID:       pid,
			Request:   &req,
			FieldMap:  fieldmap,
			ChangeMap: changemap,
			Subfields: subfields,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
		b["Fields"] = overviewfields

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}
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

		rid, err := util.GetID(c, "id")
		if err != nil {
			if err == util.ErrNoID {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Change this param and elsewhere to "fieldtype"
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

		fd, err := request.GetFieldDefinition(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !request.IsEditable(pid, &req, fd) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsFieldValueValid(req.Type, field.Type, in.Value) {
			c.Status(fiber.StatusBadRequest)
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

func CreateRequestSubfield(i *service.Interfaces) fiber.Handler {
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

		rid, err := util.GetID(c, "rid")
		if err != nil {
			if err == util.ErrNoID {
				c.Status(fiber.StatusBadRequest)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: For CreateRequestSubfield, we could probably just use the Field ID
		rfid, err := util.GetID(c, "rfid")
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

		field, err := qtx.GetRequestField(context.Background(), rfid)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		fd, err := request.GetFieldDefinition(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !request.IsEditable(pid, &req, fd) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		sfc, err := request.FieldSubfieldConfig(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !sfc.Require {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		subfields, err := qtx.ListRequestSubfieldsForField(context.Background(), field.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Log this out
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if len(subfields) >= sfc.MaxValues {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if sfc.Unique {
			for _, subfield := range subfields {
				if subfield.Value == in.Value {
					c.Status(fiber.StatusConflict)
					return nil
				}
			}
		}

		if !request.IsFieldValueValid(req.Type, field.Type, in.Value) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if err := qtx.CreateRequestSubfield(context.Background(), query.CreateRequestSubfieldParams{
			RFID:  field.ID,
			Value: in.Value,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Let this append a new element instead
		c.Append(header.HXRefresh, header.True)
		c.Status(fiber.StatusCreated)
		return nil
	}
}

func UpdateRequestSubfield(i *service.Interfaces) fiber.Handler {
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

		id, err := util.GetID(c, "id")
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

		subfield, err := qtx.GetRequestSubfield(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		field, err := qtx.GetRequestField(context.Background(), subfield.RFID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This means there's a subfield in the system without a field
				c.Status(fiber.StatusInternalServerError)
				return nil
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		req, err := qtx.GetRequest(context.Background(), field.RID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This means there's a subfield and field without a request
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		fd, err := request.GetFieldDefinition(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !request.IsEditable(pid, &req, fd) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		sfc, err := request.FieldSubfieldConfig(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !sfc.Require {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if !request.IsFieldTypeValid(req.Type, field.Type) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !request.IsFieldValueValid(req.Type, field.Type, in.Value) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if subfield.Value == in.Value {
			c.Status(fiber.StatusConflict)
			return nil
		}

		subfields, err := qtx.ListRequestSubfieldsForField(context.Background(), field.ID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Make this a utility function?
		if sfc.Unique {
			for _, subfield := range subfields {
				if subfield.Value == in.Value {
					c.Status(fiber.StatusConflict)
					return nil
				}
			}
		}

		if err := qtx.UpdateRequestSubfield(context.Background(), query.UpdateRequestSubfieldParams{
			ID:    subfield.ID,
			Value: in.Value,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Let this swap on the id
		c.Append(header.HXRefresh, header.True)
		return nil
	}
}

func DeleteRequestSubfield(i *service.Interfaces) fiber.Handler {
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

		id, err := util.GetID(c, "id")
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

		subfield, err := qtx.GetRequestSubfield(context.Background(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.Status(fiber.StatusNotFound)
				return nil
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		field, err := qtx.GetRequestField(context.Background(), subfield.RFID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This means there's a subfield in the system without a field
				c.Status(fiber.StatusInternalServerError)
				return nil
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		req, err := qtx.GetRequest(context.Background(), field.RID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: This means there's a subfield and field without a request
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		fd, err := request.GetFieldDefinition(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if !request.IsEditable(pid, &req, fd) {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		sfc, err := request.FieldSubfieldConfig(req.Type, field.Type)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		if !sfc.Require {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		subfields, err := qtx.ListRequestSubfieldsForField(context.Background(), field.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Log this out
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if len(subfields)-1 < sfc.MinValues {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err := qtx.DeleteRequestSubfield(context.Background(), subfield.ID); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err := tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Append(header.HXRefresh, header.True)
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

		rid, err := util.GetID(c, "id")
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

		if status == request.StatusFulfilled {
			// TODO: Return Forbidden, Conflict, etc depending on the error
			if err := request.Fulfill(qtx, pid, &req); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		} else {
			if err := request.UpdateStatus(qtx, request.UpdateStatusParams{
				Request: &req,
				PID:     pid,
				Status:  status,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		// TODO: Retrieve the fields and check the Change Requests
		// For any change requests, at this point they should be moved from Open to Change Request
		changes, err := qtx.ListOpenRequestChangeRequestsForRequest(context.Background(), rid)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Acceptable, this means there are no change requests
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}
		changeids := []int64{}
		for _, change := range changes {
			changeids = append(changeids, change.ID)
		}
		if err = qtx.BatchCreateRequestChangeRequest(context.Background(), changeids); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		if err = qtx.BatchDeleteOpenRequestChangeRequest(context.Background(), changeids); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Success notice?
		c.Append(header.HXRefresh, header.True)
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
			c.Status(fiber.StatusForbidden)
			return nil
		}

		rid, err := util.GetID(c, "id")
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

		change, err := qtx.GetOpenRequestChangeRequestForRequestField(context.Background(), field.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Acceptable, it means there's no change request
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		// TODO: Get this in a request utility
		var status string
		if change.ID == 0 {
			status = request.FieldStatusApproved
		} else {
			status = request.FieldStatusReviewed
		}

		if err = qtx.UpdateRequestFieldStatus(context.Background(), query.UpdateRequestFieldStatusParams{
			ID:     field.ID,
			Status: status,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if status == request.FieldStatusApproved {
			change, err := qtx.GetRequestChangeRequestByFieldID(context.Background(), field.ID)
			if err != nil {
				if err == sql.ErrNoRows {
					// TODO: Acceptable; this just means there's no change request
				} else {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
			}
			if err != sql.ErrNoRows {
				if err = qtx.CreatePastRequestChangeRequest(context.Background(), change.ID); err != nil {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}

				if err = qtx.DeleteRequestChangeRequest(context.Background(), change.ID); err != nil {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}
			}
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

		rid, err := util.GetID(c, "id")
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
			Request: &req,
			PID:     pid,
			Status:  status,
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

		rid, err := util.GetID(c, "id")
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

		// TODO: Make sure this makes sense and is tested; this is meant to prevent double-creating
		// an open change request
		_, err = qtx.GetOpenRequestChangeRequestForRequestField(context.Background(), field.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				// TODO: Acceptable, means there is no Open Change Request
			} else {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}
		if err != sql.ErrNoRows {
			c.Status(fiber.StatusForbidden)
			return nil
		}

		if err = qtx.CreateOpenRequestChangeRequest(context.Background(), query.CreateOpenRequestChangeRequestParams{
			Value: field.Value,
			RFID:  field.ID,
			PID:   pid,
			Text:  text,
		}); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if field.Status == request.FieldStatusApproved {
			if err = qtx.UpdateRequestFieldStatus(context.Background(), query.UpdateRequestFieldStatusParams{
				ID:     field.ID,
				Status: request.FieldStatusReviewed,
			}); err != nil {
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Look into returning a Boost or specific components here

		c.Append(header.HXRefresh, "true")
		c.Status(fiber.StatusCreated)
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

		id, err := util.GetID(c, "id")
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

		if err = qtx.DeleteOpenRequestChangeRequest(context.Background(), change.ID); err != nil {
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

		id, err := util.GetID(c, "id")
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
		reqs, err := qtx.ListRequestsForPlayer(context.Background(), pid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c))
		}

		// TODO: Get this into a standard API on the request package
		summaries := []request.SummaryForQueue{}
		for _, req := range reqs {
			fields, err := qtx.ListRequestFieldsForRequest(context.Background(), req.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c))
			}
			summary, err := request.NewSummaryForQueue(request.NewSummaryForQueueParams{
				Query:               qtx,
				Request:             &req,
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
		b["HasCharacterApplications"] = len(reqs) > 0
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
		reqs, err := qtx.ListRequestsByTypeAndStatus(context.Background(), query.ListRequestsByTypeAndStatusParams{
			Type: request.TypeCharacterApplication,
			// TODO: Move this to a standard list
			Statuses: []string{
				request.StatusSubmitted,
				request.StatusInReview,
			},
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
		}

		summaries := []request.SummaryForQueue{}
		for _, req := range reqs {
			if req.PID == pid {
				continue
			}

			fields, err := qtx.ListRequestFieldsForRequest(context.Background(), req.ID)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.Render(view.InternalServerError, view.Bind(c), layout.Standalone)
			}
			fieldmap := request.FieldMap(fields)
			summary, err := request.NewSummaryForQueue(request.NewSummaryForQueueParams{
				Query:               qtx,
				Request:             &req,
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
		if len(summaries) > 0 {
			b["CharacterApplicationSummaries"] = summaries
		}
		return c.Render(view.CharacterApplicationQueue, b)
	}
}
