package actor

import (
	"regexp"

	"petrichormud.com/app/internal/sanitize"
)

var ImageNameSanitizer sanitize.StringRegexSanitizer = sanitize.NewStringRegexSanitizer(regexp.MustCompile(ImageNameRegex))

func SanitizeImageName(s string) string {
	return ImageNameSanitizer.Sanitize(s)
}

var CharacterNameSanitizer sanitize.StringRegexSanitizer = sanitize.NewStringRegexSanitizer(regexp.MustCompile(CharacterNameRegex))

func SanitizeCharacterName(s string) string {
	return CharacterNameSanitizer.Sanitize(s)
}

type CharacterGenderCustomSanitizer struct{}

func (z *CharacterGenderCustomSanitizer) Sanitize(s string) string {
	if s != GenderMale && s != GenderFemale && s != GenderNonBinary {
		return GenderNonBinary
	}
	return s
}

var CharacterGenderSanitizer CharacterGenderCustomSanitizer = CharacterGenderCustomSanitizer{}

func SanitizeGender(s string) string {
	return CharacterGenderSanitizer.Sanitize(s)
}

var ShortDescriptionSanitizer sanitize.StringRegexSanitizer = sanitize.NewStringRegexSanitizer(regexp.MustCompile(ShortDescriptionRegex))

func SanitizeShortDescription(s string) string {
	return ShortDescriptionSanitizer.Sanitize(s)
}

var DescriptionSanitizer sanitize.StringRegexSanitizer = sanitize.NewStringRegexSanitizer(regexp.MustCompile(DescriptionRegex))

func SanitizeDescription(s string) string {
	return DescriptionSanitizer.Sanitize(s)
}

var CharacterBackstorySanitizer sanitize.StringRegexSanitizer = sanitize.NewStringRegexSanitizer(regexp.MustCompile(CharacterBackstoryRegex))

func SanitizeCharacterBackstory(s string) string {
	return CharacterBackstorySanitizer.Sanitize(s)
}
