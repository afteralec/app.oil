package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
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
	SummaryForQueue(p SummaryForQueueParams) SummaryForQueue
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

type ReviewDialogData struct {
	Path     string
	Variable string
}

// TODO: ReviewDialog needs consolidated and cleaned up here
type SummaryForQueue struct {
	StatusIcon   StatusIcon
	ReviewDialog ReviewDialogData
	StatusColor  string
	StatusText   string
	Title        string
	Link         string
	ReviewerText template.HTML
	ID           int64
	PID          int64
}

type SummaryForQueueParams struct {
	Content content
	Request *query.Request
}

// TODO: Error output
func (d *DefaultDefinition) SummaryForQueue(p SummaryForQueueParams) SummaryForQueue {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return SummaryForQueue{}
	}

	title := def.TitleForSummary(p.Content)

	var reviewerText template.HTML
	if p.Request.Status == StatusInReview {
		reviewerText = template.HTML("<span class=\"font-semibold\">Being reviewed by:</span>")
	} else {
		reviewerText = template.HTML("<span class=\"font-semibold\">Never reviewed</span>")
	}

	reviewDialog := ReviewDialogData{
		Path:     route.RequestStatusPath(p.Request.ID),
		Variable: fmt.Sprintf("showReviewDialogFor%s", title),
	}

	// TODO: Make this resilient to a request with an invalid status
	return SummaryForQueue{
		ID:           p.Request.ID,
		PID:          p.Request.PID,
		Title:        title,
		Link:         route.RequestPath(p.Request.ID),
		StatusIcon:   NewStatusIcon(StatusIconParams{Status: p.Request.Status, Size: "48", IncludeText: false}),
		StatusColor:  StatusColors[p.Request.Status],
		StatusText:   StatusTexts[p.Request.Status],
		ReviewerText: reviewerText,
		ReviewDialog: reviewDialog,
	}
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

// TODO: Make this a standard utility
func ContentBytes(content any) ([]byte, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

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
// TODO: Clean this up
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

func NewSummaryForQueue(p SummaryForQueueParams) (SummaryForQueue, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return SummaryForQueue{}, ErrNoDefinition
	}

	return def.SummaryForQueue(p), nil
}

func IsFieldNameValid(t, name string) bool {
	definition, ok := Definitions.Get(t)
	if !ok {
		return false
	}
	return definition.IsFieldNameValid(name)
}
