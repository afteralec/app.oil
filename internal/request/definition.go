package request

import (
	"encoding/json"
	"errors"

	"petrichormud.com/app/internal/query"
)

const errNoDefinition string = "no definition with type"

var ErrNoDefinition error = errors.New(errNoDefinition)

type Definition interface {
	Type() string
	Dialogs() Dialogs
	Fields() Fields
	IsFieldNameValid(f string) bool
	ContentBytes(q *query.Queries, rid int64) ([]byte, error)
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	SummaryTitle(content map[string]string) string
	SummaryFields(p GetSummaryFieldsParams) []SummaryField
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
	definition, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return ErrNoDefinition
	}
	if err := definition.UpdateField(q, p); err != nil {
		return err
	}
	return nil
}

func View(t, f string) string {
	definition, ok := Definitions.Get(t)
	if !ok {
		return ""
	}
	fields := definition.Fields().Map
	field := fields[f]
	return field.View
}

// TODO: Make this map a comprehensive type with methods on it?
// For example, it could have a Type field, or an IsMember method, etc
func Content(q *query.Queries, req *query.Request) (map[string]string, error) {
	var b []byte
	m := map[string]string{}

	if !IsTypeValid(req.Type) {
		return m, ErrInvalidType
	}

	definition, ok := Definitions.Get(req.Type)
	if !ok {
		return m, ErrNoDefinition
	}
	b, err := definition.ContentBytes(q, req.ID)
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]string{}, err
	}

	return m, nil
}

// TODO: Clean this up based on the Fields or new Content API
func NextIncompleteField(t string, content map[string]string) (string, bool) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", false
	}
	fields := definition.Fields()
	for i, field := range fields.List {
		value, ok := content[field.Name]
		if !ok {
			continue
		}
		if len(value) == 0 {
			return field.Name, i == len(fields.List)-1
		}
	}
	return "", false
}

// TODO: Possible error output
func GetFieldLabelAndDescription(t, f string) (string, string) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", ""
	}
	fields := definition.Fields().Map
	field := fields[f]
	return field.Label, field.Description
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
}

type GetSummaryFieldsParams struct {
	Request *query.Request
	Content map[string]string
	PID     int64
}

func SummaryFields(p GetSummaryFieldsParams) []SummaryField {
	definition, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return []SummaryField{}
	}
	return definition.SummaryFields(p)
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
