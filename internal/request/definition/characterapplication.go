package definition

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/dialog"
	"petrichormud.com/app/internal/request/field"
	"petrichormud.com/app/internal/route"
)

// TODO: Get this in a shared package
var (
	ErrMissingField      error = errors.New("a field is missing")
	ErrCurrentActorImage error = errors.New("this player already has a current actor")
)

// TODO: Create constants for FieldTypes and lift them into the Request

type fulfillerCharacterApplication struct{}

// TODO: This should return a constant
func (f *fulfillerCharacterApplication) For() string {
	return "Player"
}

// TODO: Split these individual steps out into their own functions?
func (f *fulfillerCharacterApplication) Fulfill(q *query.Queries, req *query.Request) error {
	currentcount, err := q.CountCurrentActorImagePlayerPropertiesForPlayer(context.Background(), req.PID)
	if err != nil {
		return err
	}
	if currentcount > 0 {
		return ErrCurrentActorImage
	}

	fields, err := q.ListRequestFieldsForRequest(context.Background(), req.ID)
	if err != nil {
		return err
	}
	fieldmap := field.NewMap(fields)

	// TODO: Make FieldMap a valid type and add an API to it?
	// TODO: Run validations on all of these values?
	sdescfield, ok := fieldmap[FieldCharacterApplicationShortDescription.Type]
	if !ok {
		return ErrMissingField
	}
	descfield, ok := fieldmap[FieldCharacterApplicationDescription.Type]
	if !ok {
		return ErrMissingField
	}
	genderfield, ok := fieldmap[FieldCharacterApplicationGender.Type]
	if !ok {
		return ErrMissingField
	}
	namefield, ok := fieldmap[FieldCharacterApplicationName.Type]
	if !ok {
		return ErrMissingField
	}
	var nb strings.Builder
	fmt.Fprintf(&nb, "%d-%d-%s", req.PID, req.ID, namefield.Value)
	name := nb.String()
	result, err := q.CreateActorImage(context.Background(), query.CreateActorImageParams{
		Name:             name,
		Gender:           genderfield.Value,
		ShortDescription: sdescfield.Value,
		Description:      descfield.Value,
	})
	if err != nil {
		return err
	}
	aiid, err := result.LastInsertId()
	if err != nil {
		return err
	}

	keywordsfield, ok := fieldmap[FieldCharacterApplicationKeywords.Type]
	if !ok {
		return ErrMissingField
	}
	keywordsubfields, err := q.ListRequestSubfieldsForField(context.Background(), keywordsfield.ID)
	if err != nil {
		return err
	}
	for _, keywordsubfield := range keywordsubfields {
		_, err := q.CreateActorImageKeyword(context.Background(), query.CreateActorImageKeywordParams{
			AIID:    aiid,
			Keyword: keywordsubfield.Value,
		})
		if err != nil {
			return err
		}
	}

	// TODO: Create Actor Permissions

	// TODO: Maybe add a string name to each hand?
	_, err = q.CreateActorImageHand(context.Background(), query.CreateActorImageHandParams{
		AIID: aiid,
		Hand: 1,
	})
	if err != nil {
		return err
	}
	_, err = q.CreateActorImageHand(context.Background(), query.CreateActorImageHandParams{
		AIID: aiid,
		Hand: 2,
	})
	if err != nil {
		return err
	}

	// TODO: Rename Key to Type here?
	if err := q.CreateActorImageCharacterMetadata(context.Background(), query.CreateActorImageCharacterMetadataParams{
		AIID:  aiid,
		Key:   FieldCharacterApplicationName.Type,
		Value: namefield.Value,
	}); err != nil {
		return err
	}
	backstoryfield, ok := fieldmap[FieldCharacterApplicationBackstory.Type]
	if !ok {
		return ErrMissingField
	}
	if err := q.CreateActorImageCharacterMetadata(context.Background(), query.CreateActorImageCharacterMetadataParams{
		AIID:  aiid,
		Key:   FieldCharacterApplicationBackstory.Type,
		Value: backstoryfield.Value,
	}); err != nil {
		return err
	}

	result, err = q.CreateActorImagePlayerProperties(context.Background(), query.CreateActorImagePlayerPropertiesParams{
		AIID: aiid,
		PID:  req.PID,
	})
	if err != nil {
		return err
	}
	aippid, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if err := q.SetActorImagePlayerPropertiesCurrent(context.Background(), query.SetActorImagePlayerPropertiesCurrentParams{
		ID:      aippid,
		Current: true,
	}); err != nil {
		return err
	}

	return nil
}

var FulfillerCharacterApplication fulfillerCharacterApplication = fulfillerCharacterApplication{}

type titlerCharacterApplication struct{}

func (t *titlerCharacterApplication) ForOverview(fields field.Map) string {
	var sb strings.Builder
	var name string
	field, ok := fields[FieldCharacterApplicationName.Type]
	if ok {
		name = field.Value
	} else {
		name = actor.DefaultCharacterName
	}
	fmt.Fprintf(&sb, "Character Application (%s)", name)
	return sb.String()
}

var TitlerCharacterApplication titlerCharacterApplication = titlerCharacterApplication{}

var (
	FieldCharacterApplicationName             field.Field = NewFieldCharacterApplicationName()
	FieldCharacterApplicationGender           field.Field = NewFieldCharacterApplicationGender()
	FieldCharacterApplicationShortDescription field.Field = NewFieldCharacterApplicationShortDescription()
	FieldCharacterApplicationDescription      field.Field = NewFieldCharacterApplicationDescription()
	FieldCharacterApplicationBackstory        field.Field = NewFieldCharacterApplicationBackstory()
	FieldCharacterApplicationKeywords         field.Field = NewFieldCharacterApplicationKeywords()
)

var FieldsCharacterApplication field.Group = field.NewGroup([]field.Field{
	FieldCharacterApplicationName,
	FieldCharacterApplicationGender,
	FieldCharacterApplicationShortDescription,
	FieldCharacterApplicationDescription,
	FieldCharacterApplicationBackstory,
	FieldCharacterApplicationKeywords,
})

func NewFieldCharacterApplicationName() field.Field {
	b := field.FieldBuilder()
	b.Type("name")
	b.For(field.ForPlayer)
	b.Label("Name")
	b.Description("Your character's name")
	b.Help(partial.RequestFieldHelpCharacterApplicationName)
	b.Data(partial.RequestFieldDataCharacterApplicationName)
	b.Form(partial.RequestFieldFormCharacterApplicationName)
	b.FormRenderer(new(field.DefaultRenderer))
	b.Validator(&actor.CharacterNameValidator)
	return b.Build()
}

type fieldCharacterApplicationGenderFormRenderer struct{}

func (f *fieldCharacterApplicationGenderFormRenderer) Render(e *html.Engine, field *query.RequestField, _ []query.RequestSubfield, template string) (template.HTML, error) {
	b := fiber.Map{
		"FormID":     "request-form",
		"Path":       route.RequestFieldTypePath(field.RID, field.Type),
		"FieldValue": field.Value,
	}
	b["GenderRadioGroup"] = []bind.Radio{
		{
			ID:       "edit-request-character-application-gender-non-binary",
			Name:     "value",
			Variable: "gender",
			Value:    actor.GenderNonBinary,
			Label:    "Non-Binary",
			Active:   field.Value == actor.GenderNonBinary,
		},
		{
			ID:       "edit-request-character-application-gender-female",
			Name:     "value",
			Variable: "gender",
			Value:    actor.GenderFemale,
			Label:    "Female",
			Active:   field.Value == actor.GenderFemale,
		},
		{
			ID:       "edit-request-character-application-gender-male",
			Name:     "value",
			Variable: "gender",
			Value:    actor.GenderMale,
			Label:    "Male",
			Active:   field.Value == actor.GenderMale,
		},
	}
	return partial.Render(e, partial.RenderParams{
		Template: template,
		Bind:     b,
	})
}

func NewFieldCharacterApplicationGender() field.Field {
	b := field.FieldBuilder()
	b.Type("gender")
	b.For(field.ForPlayer)
	b.Label("Gender")
	b.Description("Your character's gender determines the pronouns used by third-person descriptions in the game")
	b.Help(partial.RequestFieldHelpCharacterApplicationGender)
	b.Data(partial.RequestFieldDataCharacterApplicationGender)
	b.Form(partial.RequestFieldFormCharacterApplicationGender)
	b.FormRenderer(new(fieldCharacterApplicationGenderFormRenderer))
	b.Validator(&actor.GenderValidator)
	return b.Build()
}

func NewFieldCharacterApplicationShortDescription() field.Field {
	b := field.FieldBuilder()
	b.Type("sdesc")
	b.For(field.ForPlayer)
	b.Label("Short Description")
	b.Description("This is how your character will appear in third-person descriptions during the game")
	b.Help(partial.RequestFieldHelpCharacterApplicationShortDescription)
	b.Data(partial.RequestFieldDataCharacterApplicationShortDescription)
	b.Form(partial.RequestFieldFormCharacterApplicationShortDescription)
	b.FormRenderer(new(field.DefaultRenderer))
	b.Validator(&actor.ShortDescriptionValidator)
	return b.Build()
}

func NewFieldCharacterApplicationDescription() field.Field {
	b := field.FieldBuilder()
	b.Type("desc")
	b.For(field.ForPlayer)
	b.Label("Description")
	b.Description("This is how your character will appear when examined")
	b.Help(partial.RequestFieldHelpCharacterApplicationDescription)
	b.Data(partial.RequestFieldDataCharacterApplicationDescription)
	b.Form(partial.RequestFieldFormCharacterApplicationDescription)
	b.FormRenderer(new(field.DefaultRenderer))
	b.Validator(&actor.DescriptionValidator)
	return b.Build()
}

func NewFieldCharacterApplicationBackstory() field.Field {
	b := field.FieldBuilder()
	b.Type("backstory")
	b.For(field.ForPlayer)
	b.Label("Backstory")
	b.Description("This is your character's private backstory")
	b.Help(partial.RequestFieldHelpCharacterApplicationBackstory)
	b.Data(partial.RequestFieldDataCharacterApplicationBackstory)
	b.Form(partial.RequestFieldFormCharacterApplicationBackstory)
	b.FormRenderer(new(field.DefaultRenderer))
	b.Validator(&actor.CharacterBackstoryValidator)
	return b.Build()
}

func NewFieldCharacterApplicationKeywords() field.Field {
	b := field.FieldBuilder()
	b.Type("keywords")
	b.For(field.ForReviewer)
	b.Label("Keywords")
	b.Description("These are your character's keywords")
	b.Help(partial.RequestFieldHelpCharacterApplicationKeywords)
	b.Data(partial.RequestFieldDataCharacterApplicationKeywords)
	b.Form(partial.RequestFieldFormCharacterApplicationKeywords)
	b.FormRenderer(new(field.DefaultSubfieldRenderer))
	b.Validator(&actor.KeywordValidator)
	b.SubfieldConfig(field.NewSubfieldConfig(2, 10))
	return b.Build()
}

var DialogsCharacterApplication dialog.DefinitionGroup = dialog.DefinitionGroup{
	Submit: dialog.Definition{
		Header:     "Submit This Application?",
		Text:       "Once your character application is put in review, this cannot be undone.",
		ButtonText: "Submit This Application",
		Variable:   dialog.VariableSubmit,
		Type:       dialog.TypePrimary,
	},
	Cancel: dialog.Definition{
		Header:     "Cancel This Application?",
		Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
		ButtonText: "Cancel This Application",
		Variable:   dialog.VariableCancel,
		Type:       dialog.TypeDestructive,
	},
	PutInReview: dialog.Definition{
		Header:     "Put This Application In Review?",
		Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
		ButtonText: "I'm Ready to Review This Application",
		Variable:   dialog.VariablePutInReview,
		Type:       dialog.TypePrimary,
	},
	Approve: dialog.Definition{
		Header:     "Approve This Character Application?",
		Text:       template.HTML("Once approved, <span class=\"font-semibold\">this cannot be undone</span>. The character will go back to the player for them to create."),
		ButtonText: "Approve Character",
		Variable:   dialog.VariableApprove,
		Type:       dialog.TypePrimary,
	},
	FinishReview: dialog.Definition{
		Header:     "Finish Reviewing This Character Application?",
		Text:       template.HTML("Once you finish reviewing, <span class=\"font-semibold\">this cannot be undone</span>. It will be sent back for the player to update and re-submit. Please make sure your change requests are clear!"),
		ButtonText: "Finish Review",
		Variable:   dialog.VariableFinishReview,
		Type:       dialog.TypePrimary,
	},
	Reject: dialog.Definition{
		Header:     "Reject This Character Application?",
		Text:       template.HTML("Once rejected, this Application <span class=\"font-semibold\">cannot be re-opened</span>. Please be absolutely certain before doing this."),
		ButtonText: "Reject",
		Variable:   dialog.VariableReject,
		Type:       dialog.TypeDestructive,
	},
}
