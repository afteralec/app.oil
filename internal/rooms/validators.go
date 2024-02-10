package rooms

import "regexp"

// TODO: Get these lengths in constant
// Also, precompile these regular expressions
// Split these regexes up to provide feedback on exactly which characters don't pass validation?

func IsImageNameValid(name string) bool {
	if len(name) < 6 {
		return false
	}

	if len(name) > 50 {
		return false
	}

	re := regexp.MustCompile("[^a-z-]+")
	return !re.MatchString(name)
}

func IsTitleValid(title string) bool {
	if len(title) < 2 {
		return false
	}

	if len(title) > 150 {
		return false
	}

	re := regexp.MustCompile("[^a-zA-Z, -]+")
	return !re.MatchString(title)
}

func IsDescriptionValid(description string) bool {
	if len(description) < 50 {
		return false
	}

	if len(description) > 2000 {
		return false
	}

	re := regexp.MustCompile("[^a-zA-Z,'. -]+")
	return !re.MatchString(description)
}

func IsSizeValid(size int32) bool {
	if size < 0 {
		return false
	}

	if size > 4 {
		return false
	}

	return true
}
