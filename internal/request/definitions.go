package request

import (
	"html/template"
	"regexp"

	"petrichormud.com/app/internal/queries"
)

type DefinitionInterface interface {
	Type() string
	Dialogs() DefinitionDialogs
	Fields() []Field
	ContentBytes(q *queries.Queries, rid int64) ([]byte, error)
	UpdateField(q *queries.Queries, p UpdateFieldParams) error
	SummaryTitle(content map[string]string) string
}

type FieldsInterface interface {
	IsFieldValid(f string) bool
	IsValueValid(f, v string) bool
	NextIncomplete() string
}

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

// TODO: Change Name to Tag
// TODO: Change Min and Max to MinLen and MaxLen
type Field struct {
	Name        string
	Label       string
	Description string
	View        string
	Layout      string
	Regexes     []*regexp.Regexp
	Min         int
	Max         int
}

// TODO: Test this
func (f *Field) IsValueValid(v string) bool {
	if len(v) < f.Min {
		return false
	}

	if len(v) > f.Max {
		return false
	}

	for _, regex := range f.Regexes {
		if regex.MatchString(v) {
			return false
		}
	}

	return true
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

func MakeDefinitionMap(definitions []DefinitionInterface) map[string]DefinitionInterface {
	m := make(map[string]DefinitionInterface, len(definitions))

	for _, d := range definitions {
		m[d.Type()] = d
	}

	return m
}

func MakeTypes(definitions []DefinitionInterface) []string {
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

func MakeFieldsByType(definitions []DefinitionInterface) map[string][]Field {
	fieldsByType := make(map[string][]Field, len(definitions))

	for _, d := range definitions {
		fieldsByType[d.Type()] = d.Fields()
	}

	return fieldsByType
}

func MakeFieldNamesByType(definitions []DefinitionInterface) map[string][]string {
	fieldNamesByType := make(map[string][]string, len(definitions))

	for _, d := range definitions {
		fieldNames := MakeDefinitionFieldNames(d.Fields())
		fieldNamesByType[d.Type()] = fieldNames
	}

	return fieldNamesByType
}

func MakeFieldMapsByType(definitions []DefinitionInterface) map[string]map[string]Field {
	fieldMapsByType := make(map[string]map[string]Field, len(definitions))

	for _, d := range definitions {
		fieldMap := MakeDefinitionFieldMap(d.Fields())
		fieldMapsByType[d.Type()] = fieldMap
	}

	return fieldMapsByType
}

var (
	Definitions []DefinitionInterface = []DefinitionInterface{
		&DefinitionCharacterApplication,
	}
	DefinitionMap map[string]DefinitionInterface = MakeDefinitionMap(Definitions)
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

// TODO: Change this to just View
func GetView(t, f string) string {
	fields := FieldMapsByType[t]
	field := fields[f]
	return field.View
}

func GetFieldLabelAndDescription(t, f string) (string, string) {
	fields := FieldMapsByType[t]
	field := fields[f]
	return field.Label, field.Description
}

type GetSummaryFieldsParams struct {
	Request *queries.Request
	Content map[string]string
	PID     int64
}

// TODO: Rename to SummaryFields
func GetSummaryFields(p GetSummaryFieldsParams) []SummaryField {
	switch p.Request.Type {
	case TypeCharacterApplication:
		return DefinitionCharacterApplication.GetSummaryFields(p)
	default:
		return []SummaryField{}
	}
}
