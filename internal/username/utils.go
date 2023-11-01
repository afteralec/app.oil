package username

import (
	"regexp"
	"strings"
)

func Sanitize(u string) string {
	re := regexp.MustCompile("[^a-z0-9_-]+")
	s := re.ReplaceAllString(strings.ToLower(u), "")
	return s
}

func Validate(u string) bool {
	slen := len(u)

	if slen < 4 {
		return false
	}
	if slen > 16 {
		return false
	}

	return true
}
