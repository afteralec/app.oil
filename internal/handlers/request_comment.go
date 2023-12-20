package handlers

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/shared"
)

func CreateRequestComment(i *shared.Interfaces) fiber.Handler {
	type input struct {
		Comment string `form:"comment"`
	}
	return func(c *fiber.Ctx) error {
		r := new(input)
		if err := c.BodyParser(r); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		text := request.SanitizeComment(r.Comment)
		if !request.IsCommentValid(text) {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		pid := c.Locals("pid")
		if pid == nil {
			c.Status(fiber.StatusUnauthorized)
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
		_, ok = perms.Permissions[permissions.PlayerReviewCharacterApplicationsName]
		if !ok {
			c.Status(fiber.StatusForbidden)
			return nil
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

		fieldsByType, ok := request.FieldsByType[req.Type]
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		_, ok = fieldsByType[field]
		if !ok {
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

		cr, err := qtx.CreateRequestComment(context.Background(), queries.CreateRequestCommentParams{
			RID:   rid,
			PID:   pid.(int64),
			Text:  text,
			Field: field,
		})
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		cid, err := cr.LastInsertId()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		row, err := qtx.GetCommentWithAuthor(context.Background(), cid)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if err = tx.Commit(); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		// TODO: Move this type to the bind package
		comment := request.Comment{
			ID:         row.RequestComment.ID,
			Author:     row.Player.Username,
			Text:       row.RequestComment.Text,
			AvatarLink: "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
			CreatedAt:  row.RequestComment.CreatedAt.Unix(),
		}
		return c.Render("views/partials/request/comment/current", comment.Bind(), "")
	}
}
