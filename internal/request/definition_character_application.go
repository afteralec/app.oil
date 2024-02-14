package request

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

type CharacterApplication struct{}

// Implement the DefinitionInterface for CharacterApplication
func (app *CharacterApplication) Type() string {
	return TypeCharacterApplication
}

func (app *CharacterApplication) Dialogs() DefinitionDialogs {
	return DialogsCharacterApplication
}

func (app *CharacterApplication) Fields() []Field {
	return FieldsCharacterApplication
}

func (app *CharacterApplication) ContentBytes(q *queries.Queries, rid int64) ([]byte, error) {
	content, err := q.GetCharacterApplicationContentForRequest(context.Background(), rid)
	if err != nil {
		return []byte{}, err
	}

	b, err := json.Marshal(content)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func (app *CharacterApplication) UpdateField(q *queries.Queries, p UpdateFieldParams) error {
	switch p.Field {
	case FieldName:
		if !IsNameValid(p.Value) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentName(context.Background(), queries.UpdateCharacterApplicationContentNameParams{
			RID:  p.Request.ID,
			Name: p.Value,
		}); err != nil {
			return err
		}

	case FieldGender:
		if !IsGenderValid(p.Value) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentGender(context.Background(), queries.UpdateCharacterApplicationContentGenderParams{
			RID:    p.Request.ID,
			Gender: p.Value,
		}); err != nil {
			return err
		}
	case FieldShortDescription:
		if !IsShortDescriptionValid(p.Value) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentShortDescription(context.Background(), queries.UpdateCharacterApplicationContentShortDescriptionParams{
			RID:              p.Request.ID,
			ShortDescription: p.Value,
		}); err != nil {
			return err
		}
	case FieldDescription:
		if !IsDescriptionValid(p.Value) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentDescription(context.Background(), queries.UpdateCharacterApplicationContentDescriptionParams{
			RID:         p.Request.ID,
			Description: p.Value,
		}); err != nil {
			return err
		}
	case FieldBackstory:
		if !IsBackstoryValid(p.Value) {
			return ErrInvalidInput
		}

		if err := q.UpdateCharacterApplicationContentBackstory(context.Background(), queries.UpdateCharacterApplicationContentBackstoryParams{
			RID:       p.Request.ID,
			Backstory: p.Value,
		}); err != nil {
			return err
		}
	default:
		return ErrMalformedUpdateInput
	}

	content, err := q.GetCharacterApplicationContentForRequest(context.Background(), p.Request.ID)
	if err != nil {
		return err
	}

	ready := IsCharacterApplicationValid(&content)

	if ready && p.Request.Status == StatusIncomplete {
		if err := q.CreateHistoryForRequestStatusChange(context.Background(), queries.CreateHistoryForRequestStatusChangeParams{
			RID: p.Request.ID,
			PID: p.PID,
		}); err != nil {
			return err
		}

		if err := q.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
			ID:     p.Request.ID,
			Status: StatusReady,
		}); err != nil {
			return err
		}
	} else if !ready && p.Request.Status == StatusReady {
		if err := q.CreateHistoryForRequestStatusChange(context.Background(), queries.CreateHistoryForRequestStatusChangeParams{
			RID: p.Request.ID,
			PID: p.PID,
		}); err != nil {
			return err
		}

		if err := q.UpdateRequestStatus(context.Background(), queries.UpdateRequestStatusParams{
			ID:     p.Request.ID,
			Status: StatusIncomplete,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (app *CharacterApplication) SummaryTitle(content map[string]string) string {
	var sb strings.Builder
	titleName := constant.DefaultName
	if len(content[FieldName]) > 0 {
		titleName = content[FieldName]
	}
	fmt.Fprintf(&sb, "Character Application (%s)", titleName)
	return sb.String()
}

var FieldCharacterApplicationName Field = Field{
	Name:        "name",
	Label:       "Name",
	Description: "Your character's name",
	Min:         4,
	Max:         16,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z'-]+"),
	},
	View: view.CharacterApplicationName,
}

var FieldCharacterApplicationGender Field = Field{
	Name:        "gender",
	Label:       "Gender",
	Description: "Your character's gender determines the pronouns used by third-person descriptions in the game",
	Min:         util.MinLengthOfStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	Max:         util.MaxLengthOfStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	Regexes: []*regexp.Regexp{
		util.RegexForExactMatchStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	},
	View: view.CharacterApplicationGender,
}

var FieldCharacterApplicationShortDescription Field = Field{
	Name:        "sdesc",
	Label:       "Short Description",
	Description: "This is how your character will appear in third-person descriptions during the game",
	Min:         8,
	Max:         300,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z, -]+"),
	},
	View: view.CharacterApplicationShortDescription,
}

var FieldCharacterApplicationDescription Field = Field{
	Name:        "desc",
	Label:       "Description",
	Description: "This is how your character will appear when examined",
	Min:         32,
	Max:         2000,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z, '-.!()]+"),
	},
	View: view.CharacterApplicationDescription,
}

var FieldCharacterApplicationBackstory Field = Field{
	Name:        "backstory",
	Label:       "Backstory",
	Description: "This is your character's private backstory",
	Min:         500,
	Max:         10000,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+"),
	},
	View: view.CharacterApplicationBackstory,
}

var FieldsCharacterApplication []Field = []Field{
	FieldCharacterApplicationName,
	FieldCharacterApplicationGender,
	FieldCharacterApplicationShortDescription,
	FieldCharacterApplicationDescription,
	FieldCharacterApplicationBackstory,
}

// TODO: Get this built into the Definition
func (app *CharacterApplication) GetSummaryFields(p GetSummaryFieldsParams) []SummaryField {
	if p.Request.Type == TypeCharacterApplication {
		var basePathSB strings.Builder
		fmt.Fprintf(&basePathSB, "/requests/%d", p.Request.ID)
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

		// TODO: Build a utility for this
		allowEdit := p.Request.PID == p.PID
		if p.Request.Status != StatusIncomplete && p.Request.Status != StatusReady {
			allowEdit = false
		}

		return []SummaryField{
			{
				Label:     "Name",
				Content:   p.Content[FieldName],
				AllowEdit: allowEdit,
				Path:      namePathSB.String(),
			},
			{
				Label:     "Gender",
				Content:   p.Content[FieldGender],
				AllowEdit: allowEdit,
				Path:      genderPathSB.String(),
			},
			{
				Label:     "Short Description",
				Content:   p.Content[FieldShortDescription],
				AllowEdit: allowEdit,
				Path:      shortDescriptionPathSB.String(),
			},
			{
				Label:     "Description",
				Content:   p.Content[FieldDescription],
				AllowEdit: allowEdit,
				Path:      descriptionPathSB.String(),
			},
			{
				Label:     "Backstory",
				Content:   p.Content[FieldBackstory],
				AllowEdit: allowEdit,
				Path:      backstoryPathSB.String(),
			},
		}
	}

	return []SummaryField{}
}

var DialogsCharacterApplication DefinitionDialogs = DefinitionDialogs{
	Submit: DefinitionDialog{
		Header:     "Submit This Application?",
		Text:       "Once your character application is put in review, this cannot be undone.",
		ButtonText: "Submit This Application",
	},
	Cancel: DefinitionDialog{
		Header:     "Cancel This Application?",
		Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
		ButtonText: "Cancel This Application",
	},
	PutInReview: DefinitionDialog{
		Header:     "Put This Application In Review?",
		Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
		ButtonText: "I'm Ready to Review This Application",
	},
}

var DefinitionCharacterApplication CharacterApplication = CharacterApplication{}
