package definition

import (
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

var FieldCharacterApplicationName field.Field = NewFieldCharacterApplicationName()

var CharacterApplicationFields field.Group = field.NewGroup([]field.Field{
	FieldCharacterApplicationName,
	NewFieldCharacterApplicationGender(),
	NewFieldCharacterApplicationShortDescription(),
	NewFieldCharacterApplicationDescription(),
	NewFieldCharacterApplicationBackstory(),
})

func NewFieldCharacterApplicationName() field.Field {
	b := field.FieldBuilder()
	b.Type("name")
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

func (f *fieldCharacterApplicationGenderFormRenderer) Render(e *html.Engine, field *query.RequestField, template string) (template.HTML, error) {
	b := fiber.Map{
		"FormID":     "request-form",
		"Path":       route.RequestFieldPath(field.RID, field.Type),
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
	b.Label("Gender")
	b.Description("Your character's gender determines the pronouns used by third-person descriptions in the game")
	b.Help(partial.RequestFieldHelpCharacterApplicationGender)
	b.Data(partial.RequestFieldDataCharacterApplicationGender)
	b.Form(partial.RequestFieldFormCharacterApplicationGender)
	b.FormRenderer(new(field.DefaultRenderer))
	b.Validator(&actor.GenderValidator)
	return b.Build()
}

func NewFieldCharacterApplicationShortDescription() field.Field {
	b := field.FieldBuilder()
	b.Type("sdesc")
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
	b.Label("Backstory")
	b.Description("This is your character's private backstory")
	b.Help(partial.RequestFieldHelpCharacterApplicationBackstory)
	b.Data(partial.RequestFieldDataCharacterApplicationBackstory)
	b.Form(partial.RequestFieldFormCharacterApplicationBackstory)
	b.FormRenderer(new(field.DefaultRenderer))
	b.Validator(&actor.CharacterBackstoryValidator)
	return b.Build()
}

var DialogsCharacterApplication dialog.DefinitionGroup = dialog.DefinitionGroup{
	Submit: dialog.Definition{
		Header:     "Submit This Application?",
		Text:       "Once your character application is put in review, this cannot be undone.",
		ButtonText: "Submit This Application",
		Variable:   dialog.VariableSubmit,
	},
	Cancel: dialog.Definition{
		Header:     "Cancel This Application?",
		Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
		ButtonText: "Cancel This Application",
		Variable:   dialog.VariableCancel,
	},
	PutInReview: dialog.Definition{
		Header:     "Put This Application In Review?",
		Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
		ButtonText: "I'm Ready to Review This Application",
		Variable:   dialog.VariablePutInReview,
	},
	Approve: dialog.Definition{
		Header:     "Approve This Character Application?",
		Text:       template.HTML("Once approved, <span class=\"font-semibold\">this cannot be undone</span>. The character will go back to the player for them to create."),
		ButtonText: "Approve Character",
		Variable:   dialog.VariableApprove,
	},
	FinishReview: dialog.Definition{
		Header:     "Finish Reviewing This Character Application?",
		Text:       template.HTML("Once you finish reviewing, <span class=\"font-semibold\">this cannot be undone</span>. It will be sent back for the player to update and re-submit. Please make sure your change requests are clear!"),
		ButtonText: "Finish Review",
		Variable:   dialog.VariableFinishReview,
	},
}
