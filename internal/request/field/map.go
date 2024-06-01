package field

import "petrichormud.com/app/internal/query"

type Map = map[string]query.RequestField

func NewMap(fields []query.RequestField) Map {
	m := Map{}
	for _, field := range fields {
		m[field.Type] = field
	}
	return m
}
