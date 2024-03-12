package request

import (
	"regexp"
)

const (
	ChangeRequestTextMinLength = 10
	ChangeRequestTextMaxLength = 1000
)

func SanitizeChangeRequestText(c string) string {
	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return re.ReplaceAllString(c, "")
}

// TODO: Turn this into a Validator
func IsChangeRequestTextValid(c string) bool {
	if len(c) < ChangeRequestTextMinLength {
		return false
	}
	if len(c) > ChangeRequestTextMaxLength {
		return false
	}
	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return !re.MatchString(c)
}
