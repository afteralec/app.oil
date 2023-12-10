package character

import (
	"regexp"
	"strings"
)

func SanitizeName(u string) string {
	re := regexp.MustCompile("[^a-zA-Z'-]+")
	s := re.ReplaceAllString(strings.ToLower(u), "")
	return s
}

func IsValidName(n string) bool {
	if len(n) < 4 {
		return false
	}

	if len(n) > 16 {
		return false
	}

	return true
}
