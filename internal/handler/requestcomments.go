package handler

import (
	"context"
	"database/sql"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/interfaces"
	"petrichormud.com/app/internal/partial"
	playerpermission "petrichormud.com/app/internal/player/permission"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request"
)

func CreateRequestComment(i *interfaces.Shared) fiber.Handler {
	type input struct {
		Comment string `form:"comment"`
	}
	return func(c *fiber.Ctx) error {
		in := new(input)
		if err := c.BodyParser(in); err != nil {
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		text := request.SanitizeComment(in.Comment)
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
		perms, ok := lperms.(playerpermission.PlayerGranted)
		if !ok {
			c.Status(fiber.StatusInternalServerError)
			return nil
		}
		_, ok = perms.Permissions[playerpermission.PlayerReviewCharacterApplicationsName]
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

		fieldMapByType, ok := request.FieldMapsByType[req.Type]
		if !ok {
			c.Status(fiber.StatusBadRequest)
			return nil
		}
		_, ok = fieldMapByType[field]
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

		cr, err := qtx.CreateRequestComment(context.Background(), query.CreateRequestCommentParams{
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
			Current:        true,
			ID:             row.RequestComment.ID,
			VID:            row.RequestComment.VID,
			Author:         row.Player.Username,
			Text:           row.RequestComment.Text,
			AvatarLink:     "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
			CreatedAt:      row.RequestComment.CreatedAt.Unix(),
			ViewedByAuthor: true,
			Replies:        []request.Comment{},
		}
		return c.Render(partial.RequestCommentCurrent, comment.Bind(), "")
	}
}
