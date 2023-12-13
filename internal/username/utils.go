package username

import (
	"regexp"
	"strings"
)

const (
	MinLength = 4
	MaxLength = 16
)

func Sanitize(u string) string {
	re := regexp.MustCompile("[^a-z0-9_-]+")
	s := re.ReplaceAllString(strings.ToLower(u), "")
	return s
}

func Validate(u string) bool {
	slen := len(u)

	if slen < MinLength {
		return false
	}
	if slen > MaxLength {
		return false
	}

	return true
}
