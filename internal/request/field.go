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
}

type fieldBuilder struct {
	Field Field
}

func FieldBuilder() *fieldBuilder {
	return new(fieldBuilder)
}

func (b *fieldBuilder) Name(name string) *fieldBuilder {
	b.Field.Name = name
	return b
}

func (b *fieldBuilder) Label(label string) *fieldBuilder {
	b.Field.Label = label
	return b
}

func (b *fieldBuilder) Description(description string) *fieldBuilder {
	b.Field.Description = description
	return b
}

func (b *fieldBuilder) View(view string) *fieldBuilder {
	b.Field.View = view
	return b
}

func (b *fieldBuilder) Layout(layout string) *fieldBuilder {
	b.Field.Layout = layout
	return b
}

func (b *fieldBuilder) Updater(updater FieldUpdater) *fieldBuilder {
	b.Field.Updater = updater
	return b
}

// TODO: Separate the length validator?
func (b *fieldBuilder) Validator(validator FieldValidator) *fieldBuilder {
	// TODO: Panic if a field is attempted to be built with more than one length validator
	b.Field.Validators = append(b.Field.Validators, validator)
	return b
}

func (b *fieldBuilder) Build() Field {
	// TODO: Validate that the fields work
	return b.Field
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
