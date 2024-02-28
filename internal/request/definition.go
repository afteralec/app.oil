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
	Content(q *query.Queries, rid int64) (content, error)
	IsContentValid(c content) bool
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	TitleForSummary(c content) string
	FieldsForSummary(p FieldsForSummaryParams) ([]FieldForSummary, error)
}

type DefaultDefinition struct{}

type UpdateFieldParams struct {
	Request   *query.Request
	FieldName string
	Value     string
	PID       int64
}

func (d *DefaultDefinition) UpdateField(q *query.Queries, p UpdateFieldParams) error {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return ErrNoDefinition
	}
	fields := def.Fields()
	if err := fields.Update(q, p); err != nil {
		return err
	}

	c, err := def.Content(q, p.Request.ID)
	if err != nil {
		return err
	}
	ready := def.IsContentValid(c)
	if err := UpdateReadyStatus(q, UpdateReadyStatusParams{
		Status: p.Request.Status,
		PID:    p.PID,
		RID:    p.Request.ID,
		Ready:  ready,
	}); err != nil {
		return err
	}

	return nil
}

type FieldsForSummaryParams struct {
	Content content
	Request *query.Request
	PID     int64
}

func (d *DefaultDefinition) FieldsForSummary(p FieldsForSummaryParams) ([]FieldForSummary, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return []FieldForSummary{}, ErrNoDefinition
	}
	fields := def.Fields()
	return fields.ForSummary(p), nil
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

func FieldsForSummary(p FieldsForSummaryParams) ([]FieldForSummary, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return []FieldForSummary{}, ErrNoDefinition
	}
	return def.FieldsForSummary(p)
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

func ContentBytes(content any) ([]byte, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// TODO: Let this return the fully-qualified type
func Content(q *query.Queries, req *query.Request) (content, error) {
	if !IsTypeValid(req.Type) {
		return content{}, ErrInvalidType
	}

	definition, ok := Definitions.Get(req.Type)
	if !ok {
		return content{}, ErrNoDefinition
	}

	c, err := definition.Content(q, req.ID)
	if err != nil {
		return content{}, err
	}

	return c, nil
}

func NextIncompleteField(t string, c content) (string, bool) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", false
	}
	fields := definition.Fields()
	return fields.NextIncomplete(c)
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

func TitleForSummary(t string, c content) string {
	if !IsTypeValid(t) {
		return "Request"
	}
	definition, ok := Definitions.Get(t)
	if !ok {
		return ""
	}
	return definition.TitleForSummary(c)
}

func IsFieldNameValid(t, name string) bool {
	definition, ok := Definitions.Get(t)
	if !ok {
		return false
	}
	return definition.IsFieldNameValid(name)
}
