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

func IsValid(u string) bool {
	if len(u) < MinLength {
		return false
	}
	if len(u) > MaxLength {
		return false
	}

	return true
}
