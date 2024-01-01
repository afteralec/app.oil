package request

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/character"
	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

type SummaryField struct {
	Label          string
	Content        string
	Path           string
	ViewedByPlayer bool
}

type BindRequestPageParams struct {
	Request *queries.Request
	PID     int64
}

func BindRequestPage(b fiber.Map, p BindRequestPageParams) fiber.Map {
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

	b["StatusText"] = StatusTexts[p.Request.Status]

	b["StatusColor"] = StatusColors[p.Request.Status]

	b["ViewedByPlayer"] = p.Request.PID == p.PID
	b["ViewedByReviewer"] = p.Request.RPID == p.PID

	b["HeaderStatusIcon"] = MakeStatusIcon(p.Request.Status, 36)

	return b
}

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
		b["Header"] = "Backstory"
		b["SubHeader"] = "This is your character's private backstory."
	}

	b["RequestPath"] = routes.RequestPath(p.Request.ID)

	b["RequestFormID"] = "request-form"

	return b
}

type BindCharacterApplicationPageParams struct {
	Application    *queries.CharacterApplicationContent
	ViewedByPlayer bool
}

func BindCharacterApplicationPage(b fiber.Map, p BindCharacterApplicationPageParams) fiber.Map {
	var sb strings.Builder
	titleName := constants.DefaultName
	if len(p.Application.Name) > 0 {
		titleName = p.Application.Name
	}
	fmt.Fprintf(&sb, "Character Application (%s)", titleName)
	b["RequestTitle"] = sb.String()

	// TODO: Get this into a utility
	var basePathSB strings.Builder
	fmt.Fprintf(&basePathSB, "/requests/%d", p.Application.RID)
	basePath := basePathSB.String()

	var namePathSB strings.Builder
	fmt.Fprintf(&namePathSB, "%s/%s", basePath, FieldName)

	var genderPathSB strings.Builder
	fmt.Fprintf(&genderPathSB, "%s/%s", basePath, FieldGender)

	var shortDescriptionPathSB strings.Builder
	fmt.Fprintf(&shortDescriptionPathSB, "%s/%s", basePath, FieldShortDescription)

	var descriptionPathSB strings.Builder
	fmt.Fprintf(&descriptionPathSB, "%s/%s", basePath, FieldDescription)

	var backstoryPathSB strings.Builder
	fmt.Fprintf(&backstoryPathSB, "%s/%s", basePath, FieldBackstory)

	b["SummaryFields"] = []SummaryField{
		{
			Label:          "Name",
			Content:        p.Application.Name,
			ViewedByPlayer: p.ViewedByPlayer,
			Path:           namePathSB.String(),
		},
		{
			Label:          "Gender",
			Content:        p.Application.Gender,
			ViewedByPlayer: p.ViewedByPlayer,
			Path:           genderPathSB.String(),
		},
		{
			Label:          "Short Description",
			Content:        p.Application.ShortDescription,
			ViewedByPlayer: p.ViewedByPlayer,
			Path:           shortDescriptionPathSB.String(),
		},
		{
			Label:          "Description",
			Content:        p.Application.Description,
			ViewedByPlayer: p.ViewedByPlayer,
			Path:           descriptionPathSB.String(),
		},
		{
			Label:          "Backstory",
			Content:        p.Application.Backstory,
			ViewedByPlayer: p.ViewedByPlayer,
			Path:           backstoryPathSB.String(),
		},
	}

	return b
}

type BindDialog struct {
	Header     string
	Text       template.HTML
	ButtonText string
	Path       string
	Variable   string
}

type BindCharacterApplicationDialogsParams struct {
	Request *queries.Request
}

// TODO: Build a map that does this by request type
func BindCharacterApplicationDialogs(b fiber.Map, p BindCharacterApplicationDialogsParams) fiber.Map {
	b["CancelDialog"] = BindDialog{
		Header:     "Cancel This Application?",
		Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
		ButtonText: "Cancel This Application",
		Path:       routes.RequestPath(p.Request.ID),
		Variable:   "showCancelDialog",
	}

	b["SubmitDialog"] = BindDialog{
		Header:     "Submit This Application?",
		Text:       "Once your character application is put in review, this cannot be undone.",
		ButtonText: "Submit This Application",
		Path:       routes.RequestPath(p.Request.ID),
		Variable:   "showSubmitDialog",
	}

	b["PutInReviewDialog"] = BindDialog{
		Header:     "Put This Application In Review?",
		Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
		ButtonText: "I'm Ready to Review This Application",
		Path:       routes.RequestPath(p.Request.ID),
		Variable:   "showPutInReviewDialog",
	}

	return b
}

type BindCharacterApplicationFieldPageParams struct {
	Application *queries.CharacterApplicationContent
	Request     *queries.Request
	Field       string
}

func BindCharacterApplicationFieldPage(b fiber.Map, p BindCharacterApplicationFieldPageParams) fiber.Map {
	var sb strings.Builder
	titleName := constants.DefaultName
	if len(p.Application.Name) > 0 {
		titleName = p.Application.Name
	}
	fmt.Fprintf(&sb, "Character Application (%s)", titleName)
	b["RequestTitle"] = sb.String()

	b["Name"] = p.Application.Name
	b["Gender"] = character.SanitizeGender(p.Application.Gender)
	b["ShortDescription"] = p.Application.ShortDescription
	b["Description"] = p.Application.Description
	b["Backstory"] = p.Application.Backstory
	b["CharacterApplicationNav"] = MakeCharacterApplicationNav(p.Field, p.Application)

	if p.Field == FieldBackstory {
		// TODO: Constant; maybe a "Lastfield" declaration
		b["UpdateButtonText"] = "Finish"
	} else {
		b["UpdateButtonText"] = "Next"
	}

	// TODO: Move this field to constants?
	b["GenderNonBinary"] = character.GenderNonBinary
	b["GenderFemale"] = character.GenderFemale
	b["GenderMale"] = character.GenderMale

	b["GenderIsNonBinary"] = p.Application.Gender == character.GenderNonBinary
	b["GenderIsFemale"] = p.Application.Gender == character.GenderFemale
	b["GenderIsMale"] = p.Application.Gender == character.GenderMale

	// TODO: Get this in a declarative state too
	b["FieldName"] = FieldName
	b["FieldGender"] = FieldGender
	b["FieldShortDescription"] = FieldShortDescription
	b["FieldDescription"] = FieldDescription
	b["FieldBackstory"] = FieldBackstory

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
