package request

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/layout"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/view"
)

type CharacterApplication struct {
	DefaultDefinition
}

func (app *CharacterApplication) New(q *query.Queries, pid int64) (int64, error) {
	result, err := q.CreateRequest(context.Background(), query.CreateRequestParams{
		PID:  pid,
		Type: TypeCharacterApplication,
	})
	if err != nil {
		return 0, err
	}
	rid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err = q.CreateCharacterApplicationContent(context.Background(), rid); err != nil {
		return 0, err
	}

	if err = q.CreateCharacterApplicationContentReview(context.Background(), query.CreateCharacterApplicationContentReviewParams{
		RID:              rid,
		Name:             FieldStatusNotReviewed,
		Gender:           FieldStatusNotReviewed,
		ShortDescription: FieldStatusNotReviewed,
		Description:      FieldStatusNotReviewed,
		Backstory:        FieldStatusNotReviewed,
	}); err != nil {
		return 0, err
	}

	return rid, nil
}

func (app *CharacterApplication) Type() string {
	return TypeCharacterApplication
}

func (app *CharacterApplication) Dialogs() Dialogs {
	return DialogsCharacterApplication
}

func (app *CharacterApplication) Fields() Fields {
	return FieldsCharacterApplication
}

func (app *CharacterApplication) IsFieldNameValid(name string) bool {
	return FieldsCharacterApplication.IsFieldNameValid(name)
}

func (app *CharacterApplication) Content(q *query.Queries, rid int64) (content, error) {
	var b []byte
	m := map[string]string{}
	c, err := q.GetCharacterApplicationContentForRequest(context.Background(), rid)
	if err != nil {
		return content{}, err
	}
	b, err = ContentBytes(c)
	if err != nil {
		return content{}, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return content{}, err
	}
	return content{Inner: m}, nil
}

func (app *CharacterApplication) IsContentValid(c content) bool {
	fields := app.Fields()
	for _, field := range fields.List {
		v, ok := c.Value(field.Name)
		if !ok {
			return false
		}
		if !field.IsValueValid(v) {
			return false
		}
	}
	return true
}

func (app *CharacterApplication) ContentReview(q *query.Queries, rid int64) (contentreview, error) {
	var b []byte
	m := map[string]string{}
	c, err := q.GetCharacterApplicationContentReviewForRequest(context.Background(), rid)
	if err != nil {
		return contentreview{}, err
	}
	b, err = ContentBytes(c)
	if err != nil {
		return contentreview{}, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return contentreview{}, err
	}
	return contentreview{Inner: m}, nil
}

func (app *CharacterApplication) TitleForSummary(c content) string {
	var sb strings.Builder
	titleName := actor.DefaultCharacterName
	characterName, ok := c.Value(FieldCharacterApplicationName.Name)
	if ok {
		titleName = characterName
	}
	fmt.Fprintf(&sb, "Character Application (%s)", titleName)
	return sb.String()
}

type fieldCharacterApplicationNameUpdater struct{}

func (f *fieldCharacterApplicationNameUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentName(context.Background(), query.UpdateCharacterApplicationContentNameParams{
		RID:  p.Request.ID,
		Name: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

func (f *fieldCharacterApplicationNameUpdater) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if err := q.UpdateCharacterApplicationContentReviewName(context.Background(), query.UpdateCharacterApplicationContentReviewNameParams{
		RID:  p.Request.ID,
		Name: p.Status,
	}); err != nil {
		return err
	}

	return nil
}

type fieldCharacterApplicationNameDataRenderer struct{}

func (f *fieldCharacterApplicationNameDataRenderer) Render(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationName,
		Bind: fiber.Map{
			"FieldValue": p.FieldValue,
		},
	})
}

type fieldCharacterApplicationNameFormRenderer struct{}

func (f *fieldCharacterApplicationNameFormRenderer) Render(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationName,
		Bind: fiber.Map{
			"FormID":     FormID,
			"Path":       route.RequestFieldPath(p.Request.ID, p.FieldName),
			"FieldValue": p.FieldValue,
		},
	})
}

func NewFieldCharacterApplicationName() Field {
	b := FieldBuilder()
	b.Name("name")
	b.Label("Name")
	b.Description("Your character's name")
	b.Updater(new(fieldCharacterApplicationNameUpdater))
	b.Validator(&actor.CharacterNameValidator)
	b.View(view.RequestField)
	b.Layout(layout.Standalone)
	b.Help(partial.RequestFieldHelpCharacterApplicationName)
	b.DataRenderer(new(fieldCharacterApplicationNameDataRenderer))
	b.FormRenderer(new(fieldCharacterApplicationNameFormRenderer))

	return b.Build()
}

var FieldCharacterApplicationName Field = NewFieldCharacterApplicationName()

type fieldCharacterApplicationGenderUpdater struct{}

func (f *fieldCharacterApplicationGenderUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentGender(context.Background(), query.UpdateCharacterApplicationContentGenderParams{
		RID:    p.Request.ID,
		Gender: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

func (f *fieldCharacterApplicationGenderUpdater) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if err := q.UpdateCharacterApplicationContentReviewGender(context.Background(), query.UpdateCharacterApplicationContentReviewGenderParams{
		RID:    p.Request.ID,
		Gender: p.Status,
	}); err != nil {
		return err
	}

	return nil
}

type fieldCharacterApplicationGenderDataRenderer struct{}

func (f *fieldCharacterApplicationGenderDataRenderer) Render(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationGender,
		Bind: fiber.Map{
			"FieldValue": p.FieldValue,
		},
	})
}

type fieldCharacterApplicationGenderFormRenderer struct{}

func (f *fieldCharacterApplicationGenderFormRenderer) Render(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	b := fiber.Map{
		"FormID":     FormID,
		"Path":       route.RequestFieldPath(p.Request.ID, p.FieldName),
		"FieldValue": p.FieldValue,
	}
	gender, ok := p.Content.Value(FieldCharacterApplicationGender.Name)
	if !ok {
		return template.HTML(""), ErrInvalidInput
	}
	b["GenderRadioGroup"] = []bind.Radio{
		{
			ID:       "edit-request-character-application-gender-non-binary",
			Name:     "value",
			Variable: "gender",
			Value:    actor.GenderNonBinary,
			Label:    "Non-Binary",
			Active:   gender == actor.GenderNonBinary,
		},
		{
			ID:       "edit-request-character-application-gender-female",
			Name:     "value",
			Variable: "gender",
			Value:    actor.GenderFemale,
			Label:    "Female",
			Active:   gender == actor.GenderFemale,
		},
		{
			ID:       "edit-request-character-application-gender-male",
			Name:     "value",
			Variable: "gender",
			Value:    actor.GenderMale,
			Label:    "Male",
			Active:   gender == actor.GenderMale,
		},
	}
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationGender,
		Bind:     b,
	})
}

func NewFieldCharacterApplicationGender() Field {
	b := FieldBuilder()
	b.Name("gender")
	b.Label("Gender")
	b.Description("Your character's gender determines the pronouns used by third-person descriptions in the game")
	b.Updater(new(fieldCharacterApplicationGenderUpdater))
	b.Validator(&actor.GenderValidator)
	b.View(view.RequestField)
	b.Layout(layout.Standalone)
	b.Help(partial.RequestFieldHelpCharacterApplicationGender)
	b.DataRenderer(new(fieldCharacterApplicationGenderDataRenderer))
	b.FormRenderer(new(fieldCharacterApplicationGenderFormRenderer))

	return b.Build()
}

var FieldCharacterApplicationGender Field = NewFieldCharacterApplicationGender()

type fieldCharacterApplicationShortDescriptionUpdater struct{}

func (f *fieldCharacterApplicationShortDescriptionUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentShortDescription(context.Background(), query.UpdateCharacterApplicationContentShortDescriptionParams{
		RID:              p.Request.ID,
		ShortDescription: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

func (f *fieldCharacterApplicationShortDescriptionUpdater) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if err := q.UpdateCharacterApplicationContentReviewShortDescription(context.Background(), query.UpdateCharacterApplicationContentReviewShortDescriptionParams{
		RID:              p.Request.ID,
		ShortDescription: p.Status,
	}); err != nil {
		return err
	}

	return nil
}

type fieldCharacterApplicationShortDescriptionDataRenderer struct{}

func (f *fieldCharacterApplicationShortDescriptionDataRenderer) Render(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationShortDescription,
		Bind: fiber.Map{
			"FieldValue": p.FieldValue,
		},
	})
}

type fieldCharacterApplicationShortDescriptionFormRenderer struct{}

func (f *fieldCharacterApplicationShortDescriptionFormRenderer) Render(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationShortDescription,
		Bind: fiber.Map{
			"FormID":     FormID,
			"Path":       route.RequestFieldPath(p.Request.ID, p.FieldName),
			"FieldValue": p.FieldValue,
		},
	})
}

func NewFieldCharacterApplicationShortDescription() Field {
	b := FieldBuilder()
	b.Name("sdesc")
	b.Label("Short Description")
	b.Description("This is how your character will appear in third-person descriptions during the game")
	b.Updater(new(fieldCharacterApplicationShortDescriptionUpdater))
	b.Validator(&actor.ShortDescriptionValidator)
	b.View(view.RequestField)
	b.Layout(layout.Standalone)
	b.Help(partial.RequestFieldHelpCharacterApplicationShortDescription)
	b.DataRenderer(new(fieldCharacterApplicationShortDescriptionDataRenderer))
	b.FormRenderer(new(fieldCharacterApplicationShortDescriptionFormRenderer))

	return b.Build()
}

var FieldCharacterApplicationShortDescription Field = NewFieldCharacterApplicationShortDescription()

type fieldCharacterApplicationDescriptionUpdater struct{}

func (f *fieldCharacterApplicationDescriptionUpdater) Update(q *query.Queries, p UpdateFieldParams) error {
	if err := q.UpdateCharacterApplicationContentDescription(context.Background(), query.UpdateCharacterApplicationContentDescriptionParams{
		RID:         p.Request.ID,
		Description: p.Value,
	}); err != nil {
		return err
	}

	return nil
}

func (f *fieldCharacterApplicationDescriptionUpdater) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if err := q.UpdateCharacterApplicationContentReviewDescription(context.Background(), query.UpdateCharacterApplicationContentReviewDescriptionParams{
		RID:         p.Request.ID,
		Description: p.Status,
	}); err != nil {
		return err
	}

	return nil
}

type fieldCharacterApplicationDescriptionDataRenderer struct{}

func (f *fieldCharacterApplicationDescriptionDataRenderer) Render(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationDescription,
		Bind: fiber.Map{
			"FieldValue": p.FieldValue,
		},
	})
}

type fieldCharacterApplicationDescriptionFormRenderer struct{}

func (f *fieldCharacterApplicationDescriptionFormRenderer) Render(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationDescription,
		Bind: fiber.Map{
			"FormID":     FormID,
			"Path":       route.RequestFieldPath(p.Request.ID, p.FieldName),
			"FieldValue": p.FieldValue,
		},
	})
}

func NewFieldCharacterApplicationDescription() Field {
	b := FieldBuilder()
	b.Name("desc")
	b.Label("Description")
	b.Description("This is how your character will appear when examined")
	b.Updater(new(fieldCharacterApplicationDescriptionUpdater))
	b.Validator(&actor.DescriptionLengthValidator)
	b.Validator(&actor.DescriptionRegexValidator)
	b.View(view.RequestField)
	b.Layout(layout.Standalone)
	b.Help(partial.RequestFieldHelpCharacterApplicationDescription)
	b.DataRenderer(new(fieldCharacterApplicationDescriptionDataRenderer))
	b.FormRenderer(new(fieldCharacterApplicationDescriptionFormRenderer))

	return b.Build()
}

var FieldCharacterApplicationDescription Field = NewFieldCharacterApplicationDescription()

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

func (f *fieldCharacterApplicationBackstoryUpdater) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if err := q.UpdateCharacterApplicationContentReviewBackstory(context.Background(), query.UpdateCharacterApplicationContentReviewBackstoryParams{
		RID:       p.Request.ID,
		Backstory: p.Status,
	}); err != nil {
		return err
	}

	return nil
}

type fieldCharacterApplicationBackstoryDataRenderer struct{}

func (f *fieldCharacterApplicationBackstoryDataRenderer) Render(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationBackstory,
		Bind: fiber.Map{
			"FieldValue": p.FieldValue,
		},
	})
}

type fieldCharacterApplicationBackstoryFormRenderer struct{}

func (f *fieldCharacterApplicationBackstoryFormRenderer) Render(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: partial.RequestFieldFormCharacterApplicationBackstory,
		Bind: fiber.Map{
			"FormID":     FormID,
			"Path":       route.RequestFieldPath(p.Request.ID, p.FieldName),
			"FieldValue": p.FieldValue,
		},
	})
}

func NewFieldCharacterApplicationBackstory() Field {
	b := FieldBuilder()
	b.Name("backstory")
	b.Label("Backstory")
	b.Description("This is your character's private backstory")
	b.Updater(new(fieldCharacterApplicationBackstoryUpdater))
	b.Validator(&actor.CharacterBackstoryValidator)
	b.View(view.RequestField)
	b.Layout(layout.Standalone)
	b.Help(partial.RequestFieldHelpCharacterApplicationBackstory)
	b.DataRenderer(new(fieldCharacterApplicationBackstoryDataRenderer))
	b.FormRenderer(new(fieldCharacterApplicationBackstoryFormRenderer))

	return b.Build()
}

var FieldCharacterApplicationBackstory Field = NewFieldCharacterApplicationBackstory()

var FieldsCharacterApplication Fields = NewFields([]Field{
	FieldCharacterApplicationName,
	FieldCharacterApplicationGender,
	FieldCharacterApplicationShortDescription,
	FieldCharacterApplicationDescription,
	FieldCharacterApplicationBackstory,
})

var DialogsCharacterApplication Dialogs = Dialogs{
	Submit: Dialog{
		Header:     "Submit This Application?",
		Text:       "Once your character application is put in review, this cannot be undone.",
		ButtonText: "Submit This Application",
		Variable:   "showSubmitDialog",
	},
	Cancel: Dialog{
		Header:     "Cancel This Application?",
		Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
		ButtonText: "Cancel This Application",
		Variable:   "showCancelDialog",
	},
	PutInReview: Dialog{
		Header:     "Put This Application In Review?",
		Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
		ButtonText: "I'm Ready to Review This Application",
		Variable:   "showPutInReviewDialog",
	},
	Approve: Dialog{
		Header:     "Approve This Character Application?",
		Text:       template.HTML("Once approved, <span class=\"font-semibold\">this cannot be undone</span>. The character will go back to the player for them to create."),
		ButtonText: "Approve Character",
		Variable:   "showApproveDialog",
	},
	FinishReview: Dialog{
		Header:     "Finish Reviewing This Character Application?",
		Text:       template.HTML("Once you finish reviewing, <span class=\"font-semibold\">this cannot be undone</span>. It will be sent back for the player to update and re-submit. Please make sure your change requests are clear!"),
		ButtonText: "Finish Review",
		Variable:   "showFinishReviewDialog",
	},
}

var DefinitionCharacterApplication CharacterApplication = CharacterApplication{}
