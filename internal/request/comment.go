package request

import "regexp"

const (
	CommentMinLength = 1
	CommentMaxLength = 500
)

func SanitizeComment(c string) string {
	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return re.ReplaceAllString(c, "")
}

func IsCommentValid(c string) bool {
	if len(c) < CommentMinLength {
		return false
	}

	if len(c) > CommentMaxLength {
		return false
	}

	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return !re.MatchString(c)
}
