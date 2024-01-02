package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
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
			if err == util.ErrNoID {
				c.Status(fiber.StatusBadRequest)
				// TODO: 400 view
				return c.Render("views/500", c.Locals(constants.BindName), "views/layouts/standalone")
			}
			c.Status(fiber.StatusInternalServerError)
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
		_ = request.BindDialogs(b, request.BindDialogsParams{
			Request: &req,
		})

		_, err = request.GetContent(qtx, &req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Bind and return a view by the status of the application
		// 1. API should be like: b, view, layout, err := ResolveView(pid, req)

		// TODO:
		// 1. From here, if the status is incomplete:
		//    a. Bind the header and subheader for the field
		//    b. Bind the current content of the field, if any

		// TODO: To get a page for an incomplete request:
		// 1. Get the content
		// 2. Check the next incomplete field
		// 3. For that request, content, and field:
		//    a. Bind the appropriate values
		//    b. Return the appropriate view
		// For that field, we just need the current value and view

		// TODO: To get a page for a ready request, bind the summary values and return the summary

		// TODO:
		// 1. Get content as map via type
		// 2. Generic "Get Next Incomplete Field" function
		// 3. Update the Update handler to use /:field
		// 4. General bind for Field Pages
		// 5. Type-specific binds for Field Pages

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

			switch req.Status {
			case request.StatusIncomplete:
				field, err := request.CharacterApplicationGetNextIncompleteField(&app)
				if err != nil {
					// TODO: This means that all fields are filled out but the application is still Ready
					c.Status(fiber.StatusInternalServerError)
					return nil
				}

				// TODO: Clean this up and lift request-general details to the top
				b := c.Locals(constants.BindName).(fiber.Map)
				b = request.BindCharacterApplicationFieldPage(b, request.BindCharacterApplicationFieldPageParams{
					Application: &app,
					Request:     &req,
					Field:       field,
				})
				b = request.BindDialogs(b, request.BindDialogsParams{
					Request: &req,
				})
				b = request.BindRequestFieldPage(b, request.BindRequestFieldPageParams{
					PID:      pid,
					Field:    field,
					Request:  &req,
					Comments: []queries.ListCommentsForRequestWithAuthorRow{},
				})

				view := request.GetView(req.Type, field)

				if err = tx.Commit(); err != nil {
					c.Status(fiber.StatusInternalServerError)
					return nil
				}

				return c.Render(view, b, "layout-request-field-standalone")
			case request.StatusReady:
				b := c.Locals(constants.BindName).(fiber.Map)
				b = request.BindRequestPage(b, request.BindRequestPageParams{
					PID:     pid,
					Request: &req,
				})
				b = request.BindCharacterApplicationPage(b, request.BindCharacterApplicationPageParams{
					Application:    &app,
					ViewedByPlayer: req.PID == pid,
				})
				b = request.BindDialogs(b, request.BindDialogsParams{
					Request: &req,
				})

				return c.Render("views/requests/content/summary", b, "layout-request-summary")
			default:
				// TODO: Other views
				c.Status(fiber.StatusInternalServerError)
				return nil
			}
		} else {
			// TODO: This means that there's a request in the database with an invalid type
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
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
