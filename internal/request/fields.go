package request

import "petrichormud.com/app/internal/query"

type Fields struct {
	Map  map[string]Field
	List []Field
}

func NewFields(f []Field) Fields {
	return Fields{
		List: f,
		Map:  MakeDefinitionFieldMap(f),
	}
}

func (f *Fields) Update(q *query.Queries, p UpdateFieldParams) error {
	field, ok := f.Map[p.FieldName]
	if !ok {
		return ErrInvalidInput
	}
	return field.Update(q, p)
}
