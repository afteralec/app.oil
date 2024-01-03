package request

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/queries"
)

type BindViewedByParams struct {
	Request *queries.Request
	PID     int64
}

// TODO: Add ViewedByAdmin and maybe have ViewedByPlayer be ViewedByOwner
func BindViewedBy(b fiber.Map, p BindViewedByParams) fiber.Map {
	b["ViewedByPlayer"] = p.Request.PID == p.PID
	b["ViewedByReviewer"] = p.Request.RPID == p.PID

	return b
}

func BindComments(b fiber.Map, pid int64, vid int32, rows []queries.ListCommentsForRequestWithAuthorRow) fiber.Map {
	repliesByCID := map[int64][]Comment{}
	for _, row := range rows {
		if row.RequestComment.CID > 0 {
			reply := Comment{
				Current:        true,
				ID:             row.RequestComment.CID,
				VID:            row.RequestComment.VID,
				Author:         row.Player.Username,
				Text:           row.RequestComment.Text,
				AvatarLink:     "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
				CreatedAt:      row.RequestComment.CreatedAt.Unix(),
				ViewedByAuthor: row.RequestComment.PID == pid,
				Replies:        []Comment{},
			}

			replies, ok := repliesByCID[row.RequestComment.CID]
			if !ok {
				repliesByCID[row.RequestComment.CID] = []Comment{
					reply,
				}
			}

			repliesByCID[row.RequestComment.CID] = append(replies, reply)
		}
	}

	commentsByVID := map[int32][]Comment{}
	for _, row := range rows {
		if row.RequestComment.VID == vid {
			continue
		}

		comment := Comment{
			Current:        false,
			ID:             row.RequestComment.ID,
			VID:            row.RequestComment.VID,
			Author:         row.Player.Username,
			Text:           row.RequestComment.Text,
			AvatarLink:     "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
			CreatedAt:      row.RequestComment.CreatedAt.Unix(),
			ViewedByAuthor: row.RequestComment.PID == pid,
			Replies:        []Comment{},
		}
		replies, ok := repliesByCID[row.RequestComment.ID]
		if ok {
			comment.Replies = replies
		}

		comments, ok := commentsByVID[row.RequestComment.VID]
		if !ok {
			commentsByVID[row.RequestComment.VID] = []Comment{
				comment,
			}
		}

		commentsByVID[row.RequestComment.VID] = append(comments, comment)
	}

	current := []Comment{}
	for _, row := range rows {
		if row.RequestComment.VID == vid && row.RequestComment.CID == 0 {
			comment := Comment{
				Current:        true,
				ID:             row.RequestComment.ID,
				Author:         row.Player.Username,
				Text:           row.RequestComment.Text,
				AvatarLink:     "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
				CreatedAt:      row.RequestComment.CreatedAt.Unix(),
				ViewedByAuthor: row.RequestComment.PID == pid,
				Replies:        []Comment{},
			}

			replies, ok := repliesByCID[row.RequestComment.ID]
			if ok {
				comment.Replies = replies
			}

			current = append(current, comment)
		}
	}

	b["CurrentComments"] = current
	return b
}
