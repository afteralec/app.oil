package actors

import "regexp"

const (
	MinimumImageNameLength int = 4
	MaximumImageNameLength int = 50
)

func IsImageNameValid(name string) bool {
	if len(name) < MinimumImageNameLength {
		return false
	}

	if len(name) > MaximumImageNameLength {
		return false
	}

	re := regexp.MustCompile("[^a-z-]+")
	return !re.MatchString(name)
}
