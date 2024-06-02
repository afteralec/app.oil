package change

import (
	"regexp"

	"petrichormud.com/app/internal/sanitize"
	"petrichormud.com/app/internal/validate"
)

const (
	TextMinLength = 10
	TextMaxLength = 1000
)

// var TextRegex string = "[^a-zA-Z;, \"'\\-\\.?!()\\r\\n]+"
var TextRegex string = "[^a-zA-Z;,'\"-.!():/ ]+"

var (
	TextLengthValidator validate.StringLengthValidator       = validate.NewStringLengthValidator(TextMinLength, TextMaxLength)
	TextRegexValidator  validate.StringRegexNoMatchValidator = validate.NewStringRegexNoMatchValidator(regexp.MustCompile(TextRegex))
	TextValidator       validate.StringValidatorGroup        = validate.NewStringValidatorGroup([]validate.StringValidator{&TextLengthValidator, &TextRegexValidator})
)

var TextSanitizer sanitize.StringRegexSanitizer = sanitize.NewStringRegexSanitizer(regexp.MustCompile(TextRegex))

func SanitizeText(c string) string {
	return TextSanitizer.Sanitize(c)
}

func IsTextValid(c string) bool {
	return TextValidator.IsValid(c)
}
