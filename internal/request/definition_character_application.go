package request

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/view"
)

type CharacterApplication struct{}

func (app *CharacterApplication) Type() string {
	return TypeCharacterApplication
}

func (app *CharacterApplication) Dialogs() DefinitionDialogs {
	return DialogsCharacterApplication
}

func (app *CharacterApplication) Fields() Fields {
	return FieldsCharacterApplication
}

func (app *CharacterApplication) ContentBytes(q *query.Queries, rid int64) ([]byte, error) {
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

func (app *CharacterApplication) UpdateField(q *query.Queries, p UpdateFieldParams) error {
	fields := app.Fields()
	if err := fields.Update(q, p); err != nil {
		return err
	}

	// TODO: Get this in a utility and put this in the underlying Update calls
	content, err := q.GetCharacterApplicationContentForRequest(context.Background(), p.Request.ID)
	if err != nil {
		return err
	}

	ready := IsCharacterApplicationValid(&content)

	if ready && p.Request.Status == StatusIncomplete {
		if err := q.CreateHistoryForRequestStatusChange(context.Background(), query.CreateHistoryForRequestStatusChangeParams{
			RID: p.Request.ID,
			PID: p.PID,
		}); err != nil {
			return err
		}

		if err := q.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
			ID:     p.Request.ID,
			Status: StatusReady,
		}); err != nil {
			return err
		}
	} else if !ready && p.Request.Status == StatusReady {
		if err := q.CreateHistoryForRequestStatusChange(context.Background(), query.CreateHistoryForRequestStatusChangeParams{
			RID: p.Request.ID,
			PID: p.PID,
		}); err != nil {
			return err
		}

		if err := q.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
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

type fieldCharacterApplicationNameUpdater struct{}

var FieldCharacterApplicationNameUpdater fieldCharacterApplicationNameUpdater = fieldCharacterApplicationNameUpdater{}

func (f *fieldCharacterApplicationNameUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentName(context.Background(), query.UpdateCharacterApplicationContentNameParams{
		RID:  p.Request.ID,
		Name: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

var FieldCharacterApplicationName Field = Field{
	Name:        "name",
	Label:       "Name",
	Description: "Your character's name",
	MinLen:      4,
	MaxLen:      16,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z'-]+"),
	},
	View: view.CharacterApplicationName,
}

type fieldCharacterApplicationGenderUpdater struct{}

var FieldCharacterApplicationGenderUpdater fieldCharacterApplicationGenderUpdater = fieldCharacterApplicationGenderUpdater{}

func (f *fieldCharacterApplicationGenderUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentGender(context.Background(), query.UpdateCharacterApplicationContentGenderParams{
		RID:    p.Request.ID,
		Gender: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

var FieldCharacterApplicationGender Field = Field{
	Name:        "gender",
	Label:       "Gender",
	Description: "Your character's gender determines the pronouns used by third-person descriptions in the game",
	MinLen:      util.MinLengthOfStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	MaxLen:      util.MaxLengthOfStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	Regexes: []*regexp.Regexp{
		util.RegexForExactMatchStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	},
	View: view.CharacterApplicationGender,
}

type fieldCharacterApplicationShortDescriptionUpdater struct{}

var FieldCharacterApplicationShortDescriptionUpdater fieldCharacterApplicationShortDescriptionUpdater = fieldCharacterApplicationShortDescriptionUpdater{}

func (f *fieldCharacterApplicationShortDescriptionUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentShortDescription(context.Background(), query.UpdateCharacterApplicationContentShortDescriptionParams{
		RID:              p.Request.ID,
		ShortDescription: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

var FieldCharacterApplicationShortDescription Field = Field{
	Name:        "sdesc",
	Label:       "Short Description",
	Description: "This is how your character will appear in third-person descriptions during the game",
	MinLen:      8,
	MaxLen:      300,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z, -]+"),
	},
	View: view.CharacterApplicationShortDescription,
}

type fieldCharacterApplicationDescriptionUpdater struct{}

var FieldCharacterApplicationDescriptionUpdater fieldCharacterApplicationDescriptionUpdater = fieldCharacterApplicationDescriptionUpdater{}

func (f *fieldCharacterApplicationDescriptionUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentDescription(context.Background(), query.UpdateCharacterApplicationContentDescriptionParams{
		RID:         p.Request.ID,
		Description: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

var FieldCharacterApplicationDescription Field = Field{
	Name:        "desc",
	Label:       "Description",
	Description: "This is how your character will appear when examined",
	MinLen:      32,
	MaxLen:      2000,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z, '-.!()]+"),
	},
	View: view.CharacterApplicationDescription,
}

type fieldCharacterApplicationBackstoryUpdater struct{}

var FieldCharacterApplicationBackstoryUpdater fieldCharacterApplicationBackstoryUpdater = fieldCharacterApplicationBackstoryUpdater{}

func (f *fieldCharacterApplicationBackstoryUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentBackstory(context.Background(), query.UpdateCharacterApplicationContentBackstoryParams{
		RID:       p.Request.ID,
		Backstory: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

var FieldCharacterApplicationBackstory Field = Field{
	Name:        "backstory",
	Label:       "Backstory",
	Description: "This is your character's private backstory",
	MinLen:      500,
	MaxLen:      10000,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+"),
	},
	View: view.CharacterApplicationBackstory,
}

var FieldsCharacterApplication Fields = NewFields([]Field{
	FieldCharacterApplicationName,
	FieldCharacterApplicationGender,
	FieldCharacterApplicationShortDescription,
	FieldCharacterApplicationDescription,
	FieldCharacterApplicationBackstory,
})

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

		// TODO: Can build these from the individual Field or use the Content API
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
