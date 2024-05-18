package change

import "regexp"

const (
	TextMinLength = 10
	TextMaxLength = 1000
)

// TODO: Turn this into a Sanitizer
func SanitizeText(c string) string {
	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return re.ReplaceAllString(c, "")
}

// TODO: Turn this into a Validator
func IsTextValid(c string) bool {
	if len(c) < TextMinLength {
		return false
	}
	if len(c) > TextMaxLength {
		return false
	}
	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return !re.MatchString(c)
}
