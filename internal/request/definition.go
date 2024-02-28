package request

import (
	"encoding/json"

	"petrichormud.com/app/internal/query"
)

type Definition interface {
	Type() string
	Dialogs() Dialogs
	Fields() Fields
	ContentBytes(q *query.Queries, rid int64) ([]byte, error)
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	SummaryTitle(content map[string]string) string
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

func GetFieldLabelAndDescription(t, f string) (string, string) {
	fields := FieldMapsByType[t]
	field := fields[f]
	return field.Label, field.Description
}
