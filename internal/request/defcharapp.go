package request

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/query"
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

func NewFieldCharacterApplicationName() Field {
	updater := new(fieldCharacterApplicationNameUpdater)

	b := FieldBuilder()
	b.Name("name")
	b.Label("Name")
	b.Description("Your character's name")
	b.View(view.CharacterApplicationName)
	b.Updater(updater)
	b.Validator(&actor.CharacterNameValidator)

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

func NewFieldCharacterApplicationGender() Field {
	updater := new(fieldCharacterApplicationGenderUpdater)

	b := FieldBuilder()
	b.Name("gender")
	b.Label("Gender")
	b.Description("Your character's gender determines the pronouns used by third-person descriptions in the game")
	b.View(view.CharacterApplicationGender)
	b.Updater(updater)
	b.Validator(&actor.GenderValidator)

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

func NewFieldCharacterApplicationShortDescription() Field {
	updater := new(fieldCharacterApplicationShortDescriptionUpdater)

	b := FieldBuilder()
	b.Name("sdesc")
	b.Label("Short Description")
	b.Description("This is how your character will appear in third-person descriptions during the game")
	b.View(view.CharacterApplicationShortDescription)
	b.Updater(updater)
	b.Validator(&actor.ShortDescriptionValidator)

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

func NewFieldCharacterApplicationDescription() Field {
	updater := new(fieldCharacterApplicationDescriptionUpdater)

	b := FieldBuilder()
	b.Name("desc")
	b.Label("Description")
	b.Description("This is how your character will appear when examined")
	b.View(view.CharacterApplicationDescription)
	b.Updater(updater)
	b.Validator(&actor.DescriptionLengthValidator)
	b.Validator(&actor.DescriptionRegexValidator)

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

func NewFieldCharacterApplicationBackstory() Field {
	updater := new(fieldCharacterApplicationBackstoryUpdater)

	b := FieldBuilder()
	b.Name("backstory")
	b.Label("Backstory")
	b.Description("This is your character's private backstory")
	b.View(view.CharacterApplicationBackstory)
	b.Updater(updater)
	b.Validator(&actor.CharacterBackstoryValidator)

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
		ButtonText: "Finish Review",
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
