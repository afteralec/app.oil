package actor

import (
	"regexp"

	"petrichormud.com/app/internal/constant"
	"petrichormud.com/app/internal/util"
	"petrichormud.com/app/internal/validate"
)

const (
	MinimumImageNameLength        int = 4
	MaximumImageNameLength        int = 50
	MinimumShortDescriptionLength int = 8
	MaximumShortDescriptionLength int = 300
	MinimumDescriptionLength      int = 32
	MaximumDescriptionLength      int = 2000
)

var (
	ImageNameLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(MinimumImageNameLength, MaximumImageNameLength)
	ImageNameRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile("[^a-z-]+"))
	ImageNameValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&ImageNameLengthValidator, &ImageNameRegexValidator})
)

func IsImageNameValid(name string) bool {
	return ImageNameValidator.IsValid(name)
}

var GenderLengthValidator validate.StringLengthValidator = validate.NewStringLengthValidator(
	util.MinLengthOfStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
	util.MaxLengthOfStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
)

var GenderRegexValidator validate.StringRegexMatchValidator = validate.NewStringRegexMatchValidator(
	util.RegexForExactMatchStrings([]string{constant.GenderNonBinary, constant.GenderFemale, constant.GenderMale}),
)

var GenderValidator validate.StringValidatorGroup = validate.NewStringValidatorGroup([]validate.StringValidator{&GenderLengthValidator, &GenderRegexValidator})

func IsGenderValid(gender string) bool {
	return GenderValidator.IsValid(gender)
}

var (
	ShortDescriptionLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(MinimumShortDescriptionLength, MaximumShortDescriptionLength)
	ShortDescriptionRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile("[^a-zA-Z, -]+"))
	ShortDescriptionValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&ShortDescriptionLengthValidator, &ShortDescriptionRegexValidator})
)

func IsShortDescriptionValid(sdesc string) bool {
	return ShortDescriptionValidator.IsValid(sdesc)
}

var (
	DescriptionLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(MinimumDescriptionLength, MaximumDescriptionLength)
	DescriptionRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile("[^a-zA-Z, '-.!()]+"))
	DescriptionValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&DescriptionLengthValidator, &DescriptionRegexValidator})
)

func IsDescriptionValid(desc string) bool {
	return DescriptionValidator.IsValid(desc)
}

var (
	CharacterNameLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(4, 16)
	CharacterNameRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile("[^a-zA-Z'-]+"))
	CharacterNameValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&CharacterNameLengthValidator, &CharacterNameRegexValidator})
)

var (
	CharacterBackstoryLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(500, 10000)
	CharacterBackstoryRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+"))
	CharacterBackstoryValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&CharacterBackstoryLengthValidator, &CharacterBackstoryRegexValidator})
)
