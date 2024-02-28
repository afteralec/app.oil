package request

import (
	"errors"

	"petrichormud.com/app/internal/query"
)

var ErrInvalidInput error = errors.New("field value didn't pass validation")

type Fields struct {
	Map  map[string]Field
	List []Field
}

func NewFields(fields []Field) Fields {
	fieldsMap := map[string]Field{}
	for _, field := range fields {
		fieldsMap[field.Name] = field
	}
	return Fields{
		List: fields,
		Map:  fieldsMap,
	}
}

func (f *Fields) Update(q *query.Queries, p UpdateFieldParams) error {
	field, ok := f.Map[p.FieldName]
	if !ok {
		return ErrInvalidInput
	}
	return field.Update(q, p)
}

func (f *Fields) IsFieldNameValid(name string) bool {
	_, ok := f.Map[name]
	return ok
}

func (f *Fields) IsFieldValueValid(name, value string) bool {
	field, ok := f.Map[name]
	if !ok {
		return false
	}
	return field.IsValueValid(value)
}

func (f *Fields) ForSummary(p FieldsForSummaryParams) []FieldForSummary {
	result := []FieldForSummary{}
	for _, field := range f.List {
		result = append(result, field.ForSummary(p))
	}
	return result
}
