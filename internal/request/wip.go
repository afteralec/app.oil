package request

import (
	"encoding/json"
	"errors"

	"petrichormud.com/app/internal/query"
)

const errInvalidType string = "invalid type"

var ErrInvalidType error = errors.New(errInvalidType)

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
func GetNextIncompleteField(t string, content map[string]string) (string, bool) {
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

// TODO: Make this FieldTag?
type UpdateFieldParams struct {
	Request *query.Request
	Field   string
	Value   string
	PID     int64
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

func SummaryTitle(t string, content map[string]string) string {
	if !IsTypeValid(t) {
		return "Request"
	}

	definition := DefinitionMap[t]
	return definition.SummaryTitle(content)
}

type SummaryField struct {
	Label     string
	Content   string
	Path      string
	AllowEdit bool
}
