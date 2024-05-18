package request

import (
	"petrichormud.com/app/internal/request/change"
)

func SanitizeChangeRequestText(text string) string {
	return change.SanitizeText(text)
}

func IsChangeRequestTextValid(text string) bool {
	return change.IsTextValid(text)
}
