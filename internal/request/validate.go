package request

import (
	"petrichormud.com/app/internal/request/change"
)

// TODO: Move these into semantic files instead of storing together?

func IsFieldTypeValid(t, ft string) bool {
	fields, ok := FieldsByType[t]
	if !ok {
		return false
	}
	_, ok = fields.Map()[ft]
	return ok
}

func IsFieldValueValid(t, ft, v string) bool {
	fields, ok := FieldsByType[t]
	if !ok {
		return false
	}
	fd, ok := fields.Map()[ft]
	if !ok {
		return false
	}
	return fd.IsValid(v)
}

func IsStatusValid(status string) bool {
	_, ok := StatusTexts[status]
	return ok
}

func IsFieldStatusValid(status string) bool {
	switch status {
	case FieldStatusNotReviewed:
		return true
	case FieldStatusApproved:
		return true
	case FieldStatusReviewed:
		return true
	case FieldStatusRejected:
		return true
	default:
		return false
	}
}

func SanitizeChangeRequestText(text string) string {
	return change.SanitizeText(text)
}

func IsChangeRequestTextValid(text string) bool {
	return change.IsTextValid(text)
}
