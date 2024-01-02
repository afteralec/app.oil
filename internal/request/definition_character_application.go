package request

import (
	"html/template"
	"regexp"

	"petrichormud.com/app/internal/constants"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/views"
)

// TODO: Add Layout to these
var FieldsCharacterApplication []Field = []Field{
	{
		Name: FieldName,
		Min:  4,
		Max:  16,
		Regexes: []*regexp.Regexp{
			regexp.MustCompile("[^a-zA-Z'-]+"),
		},
		View: views.CharacterApplicationName,
	},
	{
		Name: FieldGender,
		Min:  util.MinLengthOfStrings([]string{constants.GenderNonBinary, constants.GenderFemale, constants.GenderMale}),
		Max:  util.MaxLengthOfStrings([]string{constants.GenderNonBinary, constants.GenderFemale, constants.GenderMale}),
		Regexes: []*regexp.Regexp{
			util.RegexForExactMatchStrings([]string{constants.GenderNonBinary, constants.GenderFemale, constants.GenderMale}),
		},
		View: views.CharacterApplicationGender,
	},
	{
		Name: FieldShortDescription,
		Min:  8,
		Max:  300,
		Regexes: []*regexp.Regexp{
			regexp.MustCompile("[^a-zA-Z, -]+"),
		},
		View: views.CharacterApplicationShortDescription,
	},
	{
		Name: FieldDescription,
		Min:  32,
		Max:  2000,
		Regexes: []*regexp.Regexp{
			regexp.MustCompile("[^a-zA-Z, '-.!()]+"),
		},
		View: views.CharacterApplicationDescription,
	},
	{
		Name: FieldBackstory,
		Min:  500,
		Max:  10000,
		Regexes: []*regexp.Regexp{
			regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+"),
		},
		View: views.CharacterApplicationBackstory,
	},
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

var DefinitionCharacterApplication Definition = NewDefinition(NewDefinitionParams{
	Type:    TypeCharacterApplication,
	Fields:  FieldsCharacterApplication,
	Dialogs: DialogsCharacterApplication,
})
