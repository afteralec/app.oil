package request

import (
	"context"
	"html/template"
	"log"
	"regexp"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

var FieldCharacterApplicationName Field = Field{
	Name:        "name",
	Label:       "Name",
	Description: "Your character's name",
	Min:         4,
	Max:         16,
	Regexes: []*regexp.Regexp{
		regexp.MustCompile("[^a-zA-Z'-]+"),
	},
	View: views.CharacterApplicationName,
}

var FieldCharacterApplicationGender Field = Field{
	Name:        "gender",
	Label:       "Gender",
	Description: "Your character's gender determines the pronouns used by third-person descriptions in the game",
	Min:         util.MinLengthOfStrings([]string{constants.GenderNonBinary, constants.GenderFemale, constants.GenderMale}),
	Max:         util.MaxLengthOfStrings([]string{constants.GenderNonBinary, constants.GenderFemale, constants.GenderMale}),
	Regexes: []*regexp.Regexp{
		util.RegexForExactMatchStrings([]string{constants.GenderNonBinary, constants.GenderFemale, constants.GenderMale}),
	},
	View: views.CharacterApplicationGender,
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
	View: views.CharacterApplicationShortDescription,
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
	View: views.CharacterApplicationDescription,
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
	View: views.CharacterApplicationBackstory,
}

var FieldsCharacterApplication []Field = []Field{
	FieldCharacterApplicationName,
	FieldCharacterApplicationGender,
	FieldCharacterApplicationShortDescription,
	FieldCharacterApplicationDescription,
	FieldCharacterApplicationBackstory,
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

type UtilitiesCharacterApplication struct {
	Content queries.CharacterApplicationContent
	RID     int64
}

func (u *UtilitiesCharacterApplication) LoadContent(qtx *queries.Queries) error {
	app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), u.RID)
	if err != nil {
		return err
	}
	u.Content = app
	return nil
}

func (u *UtilitiesCharacterApplication) IsFieldValueValid(f, v string) bool {
	return false
}

func (u *UtilitiesCharacterApplication) GetNextIncompleteField() string {
	for _, field := range FieldsCharacterApplication {
		// TODO: Let this take in the content
		log.Println(field.Name)
	}
	return ""
}

var DefinitionCharacterApplication Definition = NewDefinition(NewDefinitionParams{
	Type:    TypeCharacterApplication,
	Fields:  FieldsCharacterApplication,
	Dialogs: DialogsCharacterApplication,
})
