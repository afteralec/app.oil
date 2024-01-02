package request

import (
	"html/template"
	"regexp"
)

type Definition struct {
	Type         string
	Dialogs      DefinitionDialogs
	FieldMap     map[string]Field
	Fields       []Field
	FieldNameMap map[string]bool
	FieldNames   []string
}

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

// TODO: Validate this input on instantiation
type NewDefinitionParams struct {
	Type    string
	Dialogs DefinitionDialogs
	Fields  []Field
}

func NewDefinition(p NewDefinitionParams) Definition {
	fieldMap := MakeDefinitionFieldMap(p.Fields)
	fieldNames := MakeDefinitionFieldNames(p.Fields)
	fieldNameMap := MakeDefinitionFieldNameMap(p.Fields)

	d := Definition{
		Type:         p.Type,
		Dialogs:      p.Dialogs,
		Fields:       p.Fields,
		FieldMap:     fieldMap,
		FieldNames:   fieldNames,
		FieldNameMap: fieldNameMap,
	}

	return d
}

type Field struct {
	Name    string
	View    string
	Layout  string
	Regexes []*regexp.Regexp
	Min     int
	Max     int
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
		m[d.Type] = d
	}

	return m
}

func MakeTypes(definitions []Definition) []string {
	types := make([]string, len(definitions))

	for i, d := range definitions {
		types[i] = d.Type
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
		fieldsByType[d.Type] = d.Fields
	}

	return fieldsByType
}

func MakeFieldMapsByType(definitions []Definition) map[string]map[string]Field {
	fieldMapsByType := make(map[string]map[string]Field, len(definitions))

	for _, d := range definitions {
		fieldMapsByType[d.Type] = d.FieldMap
	}

	return fieldMapsByType
}

var (
	Definitions []Definition = []Definition{
		DefinitionCharacterApplication,
	}
	DefinitionMap map[string]Definition = MakeDefinitionMap(Definitions)
)

var (
	Types   []string = MakeTypes(Definitions)
	TypeMap          = MakeTypeMap(Types)
)

var (
	FieldsByType    map[string][]Field          = MakeFieldsByType(Definitions)
	FieldMapsByType map[string]map[string]Field = MakeFieldMapsByType(Definitions)
)

func GetView(t, f string) string {
	fields := FieldMapsByType[t]
	field := fields[f]
	return field.View
}
