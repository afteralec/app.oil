package actor

import (
	"regexp"

	"petrichormud.com/app/internal/validate"
)

var (
	ImageNameLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(ImageNameMinLen, ImageNameMaxLen)
	ImageNameRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile(ImageNameRegex))
	ImageNameValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&ImageNameLengthValidator, &ImageNameRegexValidator})
)

func IsImageNameValid(name string) bool {
	return ImageNameValidator.IsValid(name)
}

var (
	GenderLengthValidator validate.StringLengthValidator     = validate.NewStringLengthValidator(GenderMinLen, GenderMaxLen)
	GenderRegexValidator  validate.StringRegexMatchValidator = validate.NewStringRegexMatchValidator(regexp.MustCompile(GenderRegex))
)

var GenderValidator validate.StringValidatorGroup = validate.NewStringValidatorGroup([]validate.StringValidator{&GenderLengthValidator, &GenderRegexValidator})

func IsGenderValid(gender string) bool {
	return GenderValidator.IsValid(gender)
}

var (
	ShortDescriptionLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(ShortDescriptionMinLen, ShortDescriptionMaxLen)
	ShortDescriptionRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile(ShortDescriptionRegex))
	ShortDescriptionValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&ShortDescriptionLengthValidator, &ShortDescriptionRegexValidator})
)

func IsShortDescriptionValid(sdesc string) bool {
	return ShortDescriptionValidator.IsValid(sdesc)
}

var (
	DescriptionLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(DescriptionMinLen, DescriptionMaxLen)
	DescriptionRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile(DescriptionRegex))
	DescriptionValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&DescriptionLengthValidator, &DescriptionRegexValidator})
)

func IsDescriptionValid(desc string) bool {
	return DescriptionValidator.IsValid(desc)
}

var (
	CharacterNameLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(CharacterNameMinLen, CharacterNameMaxLen)
	CharacterNameRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile(CharacterNameRegex))
	CharacterNameValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&CharacterNameLengthValidator, &CharacterNameRegexValidator})
)

func IsCharacterNameValid(name string) bool {
	return CharacterNameValidator.IsValid(name)
}

var (
	CharacterBackstoryLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(CharacterBackstoryMinLen, CharacterBackstoryMaxLen)
	CharacterBackstoryRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile(CharacterBackstoryRegex))
	CharacterBackstoryValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&CharacterBackstoryLengthValidator, &CharacterBackstoryRegexValidator})
)
