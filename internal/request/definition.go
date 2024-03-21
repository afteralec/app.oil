package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

const errNoDefinition string = "no definition with type"

var ErrNoDefinition error = errors.New(errNoDefinition)

type Definition interface {
	New(q *query.Queries, pid int64) (int64, error)
	Type() string
	Dialogs() Dialogs
	Fields() Fields
	IsFieldNameValid(f string) bool
	Content(q *query.Queries, rid int64) (content, error)
	IsContentValid(c content) bool
	ContentReview(q *query.Queries, rid int64) (contentreview, error)
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	TitleForSummary(c content) string
	FieldsForSummary(p FieldsForSummaryParams) ([]FieldForSummary, error)
	SummaryForQueue(p SummaryForQueueParams) (SummaryForQueue, error)
	UpdateFieldStatus(q *query.Queries, p UpdateFieldStatusParams) error
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

	if err := fields.UpdateStatus(q, UpdateFieldStatusParams{
		FieldName: p.FieldName,
		Request:   p.Request,
		Status:    FieldStatusNotReviewed,
	}); err != nil {
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

type UpdateFieldStatusParams struct {
	Request   *query.Request
	FieldName string
	Status    string
	PID       int64
}

func (d *DefaultDefinition) UpdateFieldStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return ErrNoDefinition
	}
	fields := def.Fields()
	if err := fields.UpdateStatus(q, p); err != nil {
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
	PutInReviewDialog Dialog
	StatusColor       string
	StatusText        string
	Title             string
	Link              string
	AuthorUsername    string
	ReviewerText      template.HTML
	StatusIcon        StatusIcon
	ID                int64
	PID               int64
	ShowPutInReview   bool
}

type SummaryForQueueParams struct {
	Query               *query.Queries
	Content             content
	Request             *query.Request
	ReviewerPermissions *player.Permissions
	PlayerUsername      string
	ReviewerUsername    string
	PID                 int64
}

// TODO: Error output
func (d *DefaultDefinition) SummaryForQueue(p SummaryForQueueParams) (SummaryForQueue, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return SummaryForQueue{}, ErrNoDefinition
	}

	title := def.TitleForSummary(p.Content)

	reviewerText := ReviewerText(ReviewerTextParams{
		Request:          p.Request,
		ReviewerUsername: p.ReviewerUsername,
	})

	// TODO: Build a utility for this
	putInReviewDialog := def.Dialogs().PutInReview
	putInReviewDialog.Path = route.RequestStatusPath(p.Request.ID)
	putInReviewDialog.Variable = fmt.Sprintf("showReviewDialogForRequest%d", p.Request.ID)

	showPutInReview := CanBePutInReview(
		CanBePutInReviewParams{
			Request:             p.Request,
			ReviewerPermissions: p.ReviewerPermissions,
			PID:                 p.PID,
		},
	)

	// TODO: Make this resilient to a request with an invalid status
	return SummaryForQueue{
		ID:                p.Request.ID,
		PID:               p.Request.PID,
		Title:             title,
		Link:              route.RequestPath(p.Request.ID),
		StatusIcon:        NewStatusIcon(StatusIconParams{Status: p.Request.Status, IconSize: 48, IncludeText: false}),
		StatusColor:       StatusColors[p.Request.Status],
		StatusText:        StatusTexts[p.Request.Status],
		ReviewerText:      reviewerText,
		PutInReviewDialog: putInReviewDialog,
		AuthorUsername:    p.PlayerUsername,
		ShowPutInReview:   showPutInReview,
	}, nil
}

func UpdateField(q *query.Queries, p UpdateFieldParams) error {
	if !IsTypeValid(p.Request.Type) {
		return ErrInvalidType
	}
	definition, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return ErrNoDefinition
	}
	return definition.UpdateField(q, p)
}

func UpdateFieldStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	if !IsTypeValid(p.Request.Type) {
		return ErrInvalidType
	}
	definition, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return ErrNoDefinition
	}
	return definition.UpdateFieldStatus(q, p)
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

type NewParams struct {
	Type string
	PID  int64
}

func New(q *query.Queries, p NewParams) (int64, error) {
	if p.PID == 0 {
		return 0, ErrInvalidInput
	}

	if !IsTypeValid(p.Type) {
		return 0, ErrInvalidType
	}

	def, ok := Definitions.Get(p.Type)
	if !ok {
		return 0, ErrNoDefinition
	}

	return def.New(q, p.PID)
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

func ContentReview(q *query.Queries, req *query.Request) (contentreview, error) {
	if !IsTypeValid(req.Type) {
		return contentreview{}, ErrInvalidType
	}

	definition, ok := Definitions.Get(req.Type)
	if !ok {
		return contentreview{}, ErrNoDefinition
	}

	c, err := definition.ContentReview(q, req.ID)
	if err != nil {
		return contentreview{}, err
	}

	return c, nil
}

// TODO: Let this error out with no definition
func NextIncompleteField(t string, c content) (string, bool) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", false
	}
	fields := definition.Fields()
	return fields.NextIncomplete(c)
}

// TODO: Let this error out with no definition
func NextUnreviewedField(t string, cr contentreview) (string, bool) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return "", false
	}
	fields := definition.Fields()
	return fields.NextUnreviewed(cr)
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

	player, err := p.Query.GetPlayer(context.Background(), p.Request.PID)
	if err != nil {
		return SummaryForQueue{}, err
	}
	p.PlayerUsername = player.Username

	if p.Request.RPID != 0 {
		reviewer, err := p.Query.GetPlayer(context.Background(), p.Request.RPID)
		if err != nil {
			return SummaryForQueue{}, err
		}
		p.ReviewerUsername = reviewer.Username
	}

	return def.SummaryForQueue(p)
}

func IsFieldNameValid(t, name string) bool {
	definition, ok := Definitions.Get(t)
	if !ok {
		return false
	}
	return definition.IsFieldNameValid(name)
}
