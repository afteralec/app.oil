package request

import (
	"encoding/json"
	"errors"
	"html/template"
	"regexp"

	"petrichormud.com/app/internal/query"
)

const errInvalidType string = "invalid type"

var ErrInvalidType error = errors.New(errInvalidType)

type Definition interface {
	Type() string
	Dialogs() DefinitionDialogs
	Fields() Fields
	ContentBytes(q *query.Queries, rid int64) ([]byte, error)
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	SummaryTitle(content map[string]string) string
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

type Fields struct {
	Map  map[string]Field
	List []Field
}

type Field struct {
	Name        string
	Label       string
	Description string
	View        string
	Layout      string
	Updater     FieldUpdater
	Regexes     []*regexp.Regexp
	MinLen      int
	MaxLen      int
}

type FieldUpdater interface {
	Update(q *query.Queries, p UpdateFieldParams) error
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
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

func (f *Field) Update(q *query.Queries, p UpdateFieldParams) error {
	if !f.IsValueValid(p.Value) {
		return ErrInvalidInput
	}

	return f.Updater.Update(q, p)
}

// TODO: Test this
func (f *Field) IsValueValid(v string) bool {
	if len(v) < f.MinLen {
		return false
	}

	if len(v) > f.MaxLen {
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

// TODO: Make this map a comprehensive type with methods on it?
// For example, it could have a Type field, or an IsMember method, etc
func Content(q *query.Queries, req *query.Request) (map[string]string, error) {
	var b []byte
	m := map[string]string{}

	if !IsTypeValid(req.Type) {
		return m, ErrInvalidType
	}

	definition := DefinitionMap[req.Type]
	b, err := definition.ContentBytes(q, req.ID)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]string{}, err
	}

	return m, nil
}

// TODO: Key this into the Field API
func NextIncompleteField(t string, content map[string]string) (string, bool) {
	fields := FieldNamesByType[t]
	for i, field := range fields {
		value, ok := content[field]
		if !ok {
			continue
		}
		if len(value) == 0 {
			return field, i == len(fields)-1
		}
	}
	return "", false
}

type UpdateFieldParams struct {
	Request   *query.Request
	FieldName string
	Value     string
	PID       int64
}

func UpdateField(q *query.Queries, p UpdateFieldParams) error {
	if !IsTypeValid(p.Request.Type) {
		return ErrInvalidType
	}

	definition := DefinitionMap[p.Request.Type]
	if err := definition.UpdateField(q, p); err != nil {
		return err
	}

	return nil
}

func View(t, f string) string {
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
	Request *query.Request
	Content map[string]string
	PID     int64
}

func SummaryFields(p GetSummaryFieldsParams) []SummaryField {
	switch p.Request.Type {
	case TypeCharacterApplication:
		return DefinitionCharacterApplication.GetSummaryFields(p)
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
