package bind

import (
	"strconv"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
)

func CurrentView(b fiber.Map, c *fiber.Ctx) fiber.Map {
	b["HomeView"] = c.Path() == routes.Home
	b["ProfileView"] = c.Path() == routes.Profile || c.Path() == routes.Me
	b["CharactersView"] = c.Path() == routes.Characters
	b["PermissionsView"] = c.Path() == routes.PlayerPermissions
	return b
}

func CharacterApplicationContent(b fiber.Map, app *queries.CharacterApplicationContent) fiber.Map {
	b["Name"] = app.Name
	b["Gender"] = character.SanitizeGender(app.Gender)
	b["ShortDescription"] = app.ShortDescription
	b["Description"] = app.Description
	b["Backstory"] = app.Backstory
	return b
}

func CharacterApplicationPaths(b fiber.Map, req *queries.Request) fiber.Map {
	b["CharacterApplicationPath"] = routes.CharacterApplicationPath(strconv.FormatInt(req.ID, 10))
	b["CharacterApplicationNamePath"] = routes.CharacterApplicationNamePath(strconv.FormatInt(req.ID, 10))
	b["CharacterApplicationGenderPath"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(req.ID, 10))
	b["CharacterApplicationShortDescriptionPath"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(req.ID, 10))
	b["CharacterApplicationDescriptionPath"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(req.ID, 10))
	b["CharacterApplicationBackstoryPath"] = routes.CharacterApplicationBackstoryPath(strconv.FormatInt(req.ID, 10))
	b["SubmitCharacterApplicationPath"] = routes.SubmitCharacterApplicationPath(strconv.FormatInt(req.ID, 10))
	return b
}

func RequestCommentPaths(b fiber.Map, req *queries.Request, field string) fiber.Map {
	b["CreateRequestCommentPath"] = routes.CreateRequestCommentPath(strconv.FormatInt(req.ID, 10), field)
	return b
}

func RequestComments(b fiber.Map, pid int64, vid int32, rows []queries.ListCommentsForRequestWithAuthorRow) fiber.Map {
	repliesByCID := map[int64][]request.Comment{}
	for _, row := range rows {
		if row.RequestComment.CID > 0 {
			reply := request.Comment{
				ID:             row.RequestComment.CID,
				Author:         row.Player.Username,
				Text:           row.RequestComment.Text,
				AvatarLink:     "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
				CreatedAt:      row.RequestComment.CreatedAt.Unix(),
				ViewedByAuthor: row.RequestComment.PID == pid,
				Replies:        []request.Comment{},
			}

			replies, ok := repliesByCID[row.RequestComment.CID]
			if !ok {
				repliesByCID[row.RequestComment.CID] = []request.Comment{
					reply,
				}
			}

			repliesByCID[row.RequestComment.CID] = append(replies, reply)
		}
	}

	current := []request.Comment{}
	for _, row := range rows {
		if row.RequestComment.VID == vid && row.RequestComment.CID == 0 {
			comment := request.Comment{
				ID:             row.RequestComment.ID,
				Author:         row.Player.Username,
				Text:           row.RequestComment.Text,
				AvatarLink:     "https://gravatar.com/avatar/205e460b479e2e5b48aec07710c08d50.jpeg?f=y&r=m&s=256&d=retro",
				CreatedAt:      row.RequestComment.CreatedAt.Unix(),
				ViewedByAuthor: row.RequestComment.PID == pid,
				Replies:        []request.Comment{},
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

func CharacterApplicationNav(b fiber.Map, app *queries.CharacterApplicationContent, field string) fiber.Map {
	b["CharacterApplicationNav"] = character.MakeApplicationNav(field, app)
	return b
}

func CharacterApplicationHeaderStatusIcon(b fiber.Map, req *queries.Request) fiber.Map {
	b["HeaderStatusIcon"] = request.MakeStatusIcon(req.Status, 36)
	return b
}

func CharacterApplicationGender(b fiber.Map, app *queries.CharacterApplicationContent) fiber.Map {
	b["GenderNonBinary"] = character.GenderNonBinary
	b["GenderFemale"] = character.GenderFemale
	b["GenderMale"] = character.GenderMale
	b["GenderIsNonBinary"] = app.Gender == character.GenderNonBinary
	b["GenderIsFemale"] = app.Gender == character.GenderFemale
	b["GenderIsMale"] = app.Gender == character.GenderMale
	return b
}

func RequestStatus(b fiber.Map, req *queries.Request) fiber.Map {
	b["StatusIncomplete"] = req.Status == request.StatusIncomplete
	b["StatusReady"] = req.Status == request.StatusReady
	b["StatusSubmitted"] = req.Status == request.StatusSubmitted
	b["StatusInReview"] = req.Status == request.StatusInReview
	b["StatusApproved"] = req.Status == request.StatusApproved
	b["StatusReviewed"] = req.Status == request.StatusReviewed
	b["StatusRejected"] = req.Status == request.StatusRejected
	b["StatusArchived"] = req.Status == request.StatusArchived
	b["StatusCanceled"] = req.Status == request.StatusCanceled
	return b
}

func RequestViewedBy(b fiber.Map, req *queries.Request, pid int64) fiber.Map {
	b["ViewedByPlayer"] = req.PID == pid
	b["ViewedByReviewer"] = req.RPID == pid
	return b
}
