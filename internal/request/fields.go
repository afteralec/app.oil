package request

import "petrichormud.com/app/internal/query"

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

func (f *Fields) IsFieldValid(name string) bool {
	_, ok := f.Map[name]
	return ok
}
