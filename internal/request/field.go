package request

import (
	"regexp"

	"petrichormud.com/app/internal/query"
)

// TODO: Add custom validators here
type Field struct {
	Name        string
	Label       string
	Description string
	View        string
	Layout      string
	Updater     FieldUpdater
	Validators  []FieldValidator
	// Regexes     []*regexp.Regexp
	// MinLen      int
	// MaxLen      int
}

type FieldBuilder struct {
	Name        string
	Label       string
	Description string
	View        string
	Layout      string
	Updater     FieldUpdater
	Validators  []FieldValidator
}

type FieldValidator interface {
	IsValid(v string) bool
}

type FieldUpdater interface {
	Update(q *query.Queries, p UpdateFieldParams) error
}

func (f *Field) IsValueValid(v string) bool {
	for _, validator := range f.Validators {
		if !validator.IsValid(v) {
			return false
		}
	}

	return true
}

func (f *Field) Update(q *query.Queries, p UpdateFieldParams) error {
	if !f.IsValueValid(p.Value) {
		return ErrInvalidInput
	}
	return f.Updater.Update(q, p)
}

type FieldLengthValidator struct {
	MinLen int
	MaxLen int
}

func NewFieldLengthValidator(min, max int) FieldLengthValidator {
	return FieldLengthValidator{
		MinLen: min,
		MaxLen: max,
	}
}

func (f *FieldLengthValidator) IsValid(v string) bool {
	if len(v) < f.MinLen {
		return false
	}

	if len(v) > f.MaxLen {
		return false
	}

	return true
}

type FieldRegexMatchValidator struct {
	Regex *regexp.Regexp
}

func NewFieldRegexMatchValidator(regex *regexp.Regexp) FieldRegexMatchValidator {
	return FieldRegexMatchValidator{
		Regex: regex,
	}
}

func (f *FieldRegexMatchValidator) IsValid(v string) bool {
	return f.Regex.MatchString(v)
}

type FieldRegexNoMatchValidator struct {
	Regex *regexp.Regexp
}

func NewFieldRegexNoMatchValidator(regex *regexp.Regexp) FieldRegexNoMatchValidator {
	return FieldRegexNoMatchValidator{
		Regex: regex,
	}
}

func (f *FieldRegexNoMatchValidator) IsValid(v string) bool {
	return !f.Regex.MatchString(v)
}
