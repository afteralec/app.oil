package request

import (
	"errors"

	"petrichormud.com/app/internal/query"
)

const errInvalidType string = "invalid type"

var ErrInvalidType error = errors.New(errInvalidType)

// TODO: Add API to a Fields struct that can take in a field and value and return if it's valid
// Have the Fields struct be in charge of the list of fields and the map of fields by name

type definitions struct {
	Map  map[string]Definition
	List []Definition
}

func NewDefinitions(defs []Definition) definitions {
	m := make(map[string]Definition, len(defs))
	for _, d := range defs {
		m[d.Type()] = d
	}
	return definitions{
		Map:  m,
		List: defs,
	}
}

func (d *definitions) Get(t string) (Definition, bool) {
	definition, ok := d.Map[t]
	if !ok {
		return nil, false
	}
	return definition, true
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
}

var Definitions definitions = NewDefinitions([]Definition{
	&DefinitionCharacterApplication,
})

type GetSummaryFieldsParams struct {
	Request *query.Request
	Content map[string]string
	PID     int64
}

func SummaryFields(p GetSummaryFieldsParams) []SummaryField {
	switch p.Request.Type {
	case TypeCharacterApplication:
		return DefinitionCharacterApplication.SummaryFields(p)
	default:
		return []SummaryField{}
	}
}

func SummaryTitle(t string, content map[string]string) string {
	if !IsTypeValid(t) {
		return "Request"
	}

	definition, ok := Definitions.Get(t)
	if !ok {
		return ""
	}
	return definition.SummaryTitle(content)
}
