package sanitize

import "regexp"

type StringSanitizer interface {
	Sanitize(s string) string
}

type StringRegexSanitizer struct {
	Regex *regexp.Regexp
}

func NewStringRegexSanitizer(regex *regexp.Regexp) StringRegexSanitizer {
	return StringRegexSanitizer{
		Regex: regex,
	}
}

func (z *StringRegexSanitizer) Sanitize(s string) string {
	return z.Regex.ReplaceAllString(s, "")
}
