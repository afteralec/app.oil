package request

import (
	"errors"
)

const errInvalidType string = "invalid type"

var ErrInvalidType error = errors.New(errInvalidType)

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

var Definitions definitions = NewDefinitions([]Definition{
	&DefinitionCharacterApplication,
})
