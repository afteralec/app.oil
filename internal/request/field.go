package request

import (
	"fmt"
	"strings"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/validate"
)

type Field struct {
	Updater     FieldUpdater
	Validator   validate.StringValidator
	Name        string
	Label       string
	Description string
	View        string
	Layout      string
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

func (b *fieldBuilder) Validator(validator validate.StringValidator) *fieldBuilder {
	b.Field.Validator = validator
	return b
}

func (b *fieldBuilder) Build() Field {
	// TODO: Validate that the fields work
	return b.Field
}

type FieldUpdater interface {
	Update(q *query.Queries, p UpdateFieldParams) error
	UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error
}

func (f *Field) Update(q *query.Queries, p UpdateFieldParams) error {
	if !f.IsValueValid(p.Value) {
		return ErrInvalidInput
	}
	return f.Updater.Update(q, p)
}

func (f *Field) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if !IsFieldStatusValid(p.Status) {
		return ErrInvalidInput
	}

	return f.Updater.UpdateStatus(q, p)
}

type FieldForSummary struct {
	Label     string
	Value     string
	Path      string
	AllowEdit bool
}

// TODO: Error output?
func (f *Field) ForSummary(p FieldsForSummaryParams) FieldForSummary {
	v, ok := p.Content.Value(f.Name)
	if !ok {
		v = ""
	}

	var basePathSB strings.Builder
	fmt.Fprintf(&basePathSB, "/requests/%d", p.Request.ID)
	basePath := basePathSB.String()
	var pathSB strings.Builder
	fmt.Fprintf(&pathSB, "%s/%s", basePath, f.Name)

	// TODO: Build a utility for this
	allowEdit := p.Request.PID == p.PID
	if p.Request.Status != StatusIncomplete && p.Request.Status != StatusReady {
		allowEdit = false
	}

	return FieldForSummary{
		Label:     f.Label,
		Value:     v,
		Path:      pathSB.String(),
		AllowEdit: allowEdit,
	}
}

func (f *Field) IsValueValid(v string) bool {
	return f.Validator.IsValid(v)
}
