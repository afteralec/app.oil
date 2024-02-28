package request

import (
	"errors"
	"html/template"

	"petrichormud.com/app/internal/query"
)

const errInvalidType string = "invalid type"

var ErrInvalidType error = errors.New(errInvalidType)

// TODO: Add API to a Fields struct that can take in a field and value and return if it's valid
// Have the Fields struct be in charge of the list of fields and the map of fields by name

type DefinitionDialog struct {
	Header     string
	ButtonText string
	Text       template.HTML
}

type DefinitionDialogs struct {
	Submit      DefinitionDialog
	Cancel      DefinitionDialog
	PutInReview DefinitionDialog
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
}

func MakeDefinitionFieldMap(fields []Field) map[string]Field {
	fieldMap := map[string]Field{}

	for _, field := range fields {
		fieldMap[field.Name] = field
	}

	return fieldMap
}

func MakeDefinitionFieldNames(fields []Field) []string {
	fieldNames := make([]string, len(fields))

	for i, field := range fields {
		fieldNames[i] = field.Name
	}

	return fieldNames
}

func MakeDefinitionFieldNameMap(fields []Field) map[string]bool {
	fieldNameMap := make(map[string]bool, len(fields))

	for _, field := range fields {
		fieldNameMap[field.Name] = true
	}

	return fieldNameMap
}

func MakeDefinitionMap(definitions []Definition) map[string]Definition {
	m := make(map[string]Definition, len(definitions))

	for _, d := range definitions {
		m[d.Type()] = d
	}

	return m
}

func MakeTypes(definitions []Definition) []string {
	types := make([]string, len(definitions))

	for i, d := range definitions {
		types[i] = d.Type()
	}

	return types
}

func MakeTypeMap(types []string) map[string]bool {
	m := make(map[string]bool, len(types))

	for _, t := range types {
		m[t] = true
	}

	return m
}

func MakeFieldsByType(definitions []Definition) map[string][]Field {
	fieldsByType := make(map[string][]Field, len(definitions))

	for _, d := range definitions {
		fieldsByType[d.Type()] = d.Fields().List
	}

	return fieldsByType
}

func MakeFieldNamesByType(definitions []Definition) map[string][]string {
	fieldNamesByType := make(map[string][]string, len(definitions))

	for _, d := range definitions {
		fieldNames := MakeDefinitionFieldNames(d.Fields().List)
		fieldNamesByType[d.Type()] = fieldNames
	}

	return fieldNamesByType
}

func MakeFieldMapsByType(definitions []Definition) map[string]map[string]Field {
	fieldMapsByType := make(map[string]map[string]Field, len(definitions))

	for _, d := range definitions {
		fieldMap := MakeDefinitionFieldMap(d.Fields().List)
		fieldMapsByType[d.Type()] = fieldMap
	}

	return fieldMapsByType
}

var (
	Definitions []Definition = []Definition{
		&DefinitionCharacterApplication,
	}
	DefinitionMap map[string]Definition = MakeDefinitionMap(Definitions)
)

var (
	Types   []string = MakeTypes(Definitions)
	TypeMap          = MakeTypeMap(Types)
)

var (
	FieldsByType     map[string][]Field          = MakeFieldsByType(Definitions)
	FieldNamesByType map[string][]string         = MakeFieldNamesByType(Definitions)
	FieldMapsByType  map[string]map[string]Field = MakeFieldMapsByType(Definitions)
)

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

	definition := DefinitionMap[t]
	return definition.SummaryTitle(content)
}
