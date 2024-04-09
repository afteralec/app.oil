package request

import (
	"petrichormud.com/app/internal/request/definition"
	"petrichormud.com/app/internal/request/field"
)

type Titler interface {
	ForOverview(fields field.Map) string
}

var TitlersByType map[string]Titler = map[string]Titler{
	TypeCharacterApplication: &definition.TitlerCharacterApplication,
}

// TODO: Error output?
func Title(t string, fields field.Map) string {
	if !IsTypeValid(t) {
		return "Request"
	}
	titler, ok := TitlersByType[t]
	if !ok {
		return "Request"
	}
	return titler.ForOverview(fields)
}
