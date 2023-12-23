package request

import (
	"fmt"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

type BindRequestFieldPageParams struct {
	Field    string
	Request  *queries.Request
	Comments []queries.ListCommentsForRequestWithAuthorRow
	PID      int64
}

func BindRequestFieldPage(b fiber.Map, p BindRequestFieldPageParams) fiber.Map {
	// TODO: Turn this into iterations over the master list of request statuses
	// Same theme as below, declarative now
	b["StatusIncomplete"] = StatusIncomplete
	b["StatusReady"] = StatusReady
	b["StatusSubmitted"] = StatusSubmitted
	b["StatusInReview"] = StatusInReview
	b["StatusApproved"] = StatusApproved
	b["StatusReviewed"] = StatusReviewed
	b["StatusRejected"] = StatusRejected
	b["StatusArchived"] = StatusArchived
	b["StatusCanceled"] = StatusCanceled

	b["StatusIsIncomplete"] = p.Request.Status == StatusIncomplete
	b["StatusIsReady"] = p.Request.Status == StatusReady
	b["StatusIsSubmitted"] = p.Request.Status == StatusSubmitted
	b["StatusIsInReview"] = p.Request.Status == StatusInReview
	b["StatusIsApproved"] = p.Request.Status == StatusApproved
	b["StatusIsReviewed"] = p.Request.Status == StatusReviewed
	b["StatusIsRejected"] = p.Request.Status == StatusRejected
	b["StatusIsArchived"] = p.Request.Status == StatusArchived
	b["StatusIsCanceled"] = p.Request.Status == StatusCanceled

	b["ViewedByPlayer"] = p.Request.PID == p.PID
	b["ViewedByReviewer"] = p.Request.RPID == p.PID

	b["CreateRequestCommentPath"] = routes.CreateRequestCommentPath(strconv.FormatInt(p.Request.ID, 10), p.Field)

	// TODO: See if this can also be a Bind function
	b["HeaderStatusIcon"] = MakeStatusIcon(p.Request.Status, 36)

	b = BindComments(b, p.Request.PID, p.Request.VID, p.Comments)

	// TODO: Fix this with a map or iteration
	// Overall, this and NextLink/BackLink should be a declarative setup
	switch p.Field {
	case FieldName:
		b["Header"] = "Name"
		b["SubHeader"] = "Your character's name"
	case FieldGender:
		b["Header"] = "Gender"
		b["SubHeader"] = "Your gender determines the pronouns used by third-person descriptions in the game."
	case FieldShortDescription:
		b["Header"] = "Short Description"
		b["SubHeader"] = "This is how your character will appear in third-person descriptions during the game."
	case FieldDescription:
		b["Header"] = "Description"
		b["SubHeader"] = "This is how your character will appear when examined."
	case FieldBackstory:
		b["Header"] = "Description"
		b["SubHeader"] = "  This is your character's private backstory."
	}

	return b
}

func BindCharacterApplicationFieldPage(b fiber.Map, app *queries.CharacterApplicationContent, field string) fiber.Map {
	// TODO: Get this "Unnamed" into a constant
	var sb strings.Builder
	titleName := "Unnamed"
	if len(app.Name) > 0 {
		titleName = app.Name
	}
	fmt.Fprintf(&sb, "Character Application (%s)", titleName)
	b["RequestTitle"] = sb.String()

	b["Name"] = app.Name
	b["Gender"] = character.SanitizeGender(app.Gender)
	b["ShortDescription"] = app.ShortDescription
	b["Description"] = app.Description
	b["Backstory"] = app.Backstory
	b["CharacterApplicationNav"] = MakeCharacterApplicationNav(field, app)

	switch field {
	case FieldName:
		b["NextLink"] = routes.RequestFieldPath(app.RID, FieldGender)
	case FieldGender:
		b["BackLink"] = routes.RequestFieldPath(app.RID, FieldName)
		b["NextLink"] = routes.RequestFieldPath(app.RID, FieldShortDescription)
	case FieldShortDescription:
		b["BackLink"] = routes.RequestFieldPath(app.RID, FieldGender)
		b["NextLink"] = routes.RequestFieldPath(app.RID, FieldDescription)
	case FieldDescription:
		b["BackLink"] = routes.RequestFieldPath(app.RID, FieldShortDescription)
		b["NextLink"] = routes.RequestFieldPath(app.RID, FieldBackstory)
	case FieldBackstory:
		b["BackLink"] = routes.RequestFieldPath(app.RID, FieldDescription)
	}

	// TODO: Let's clean up these links
	b["CharacterApplicationPath"] = routes.CharacterApplicationPath(strconv.FormatInt(app.RID, 10))
	b["CharacterApplicationNamePath"] = routes.CharacterApplicationNamePath(strconv.FormatInt(app.RID, 10))
	b["CharacterApplicationGenderPath"] = routes.CharacterApplicationGenderPath(strconv.FormatInt(app.RID, 10))
	b["CharacterApplicationShortDescriptionPath"] = routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(app.RID, 10))
	b["CharacterApplicationDescriptionPath"] = routes.CharacterApplicationDescriptionPath(strconv.FormatInt(app.RID, 10))
	b["CharacterApplicationBackstoryPath"] = routes.CharacterApplicationBackstoryPath(strconv.FormatInt(app.RID, 10))
	b["SubmitCharacterApplicationPath"] = routes.SubmitCharacterApplicationPath(strconv.FormatInt(app.RID, 10))

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

type RequestNav struct {
	Link    string
	Label   string
	Current bool
	Ready   bool
}

func MakeCharacterApplicationNav(current string, app *queries.CharacterApplicationContent) []RequestNav {
	result := []RequestNav{}

	result = append(result, RequestNav{
		Label:   "Name",
		Link:    routes.RequestFieldPath(app.RID, FieldName),
		Current: current == FieldName,
		Ready:   IsNameValid(app.Name),
	})

	result = append(result, RequestNav{
		Label:   "Gender",
		Link:    routes.RequestFieldPath(app.RID, FieldGender),
		Current: current == "gender",
		Ready:   IsGenderValid(app.Gender),
	})

	result = append(result, RequestNav{
		Label:   "Short Description",
		Link:    routes.RequestFieldPath(app.RID, FieldShortDescription),
		Current: current == "sdesc",
		Ready:   IsShortDescriptionValid(app.ShortDescription),
	})

	result = append(result, RequestNav{
		Label:   "Description",
		Link:    routes.RequestFieldPath(app.RID, FieldDescription),
		Current: current == "description",
		Ready:   IsDescriptionValid(app.Description),
	})

	result = append(result, RequestNav{
		Label:   "Backstory",
		Link:    routes.RequestFieldPath(app.RID, FieldBackstory),
		Current: current == "backstory",
		Ready:   IsBackstoryValid(app.Backstory),
	})

	return result
}
