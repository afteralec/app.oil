package request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/definition"
	"petrichormud.com/app/internal/request/field"
	"petrichormud.com/app/internal/route"
)

const errNoDefinition string = "no definition with type"

var ErrNoDefinition error = errors.New(errNoDefinition)

type Definition interface {
	// Can be broken out
	New(q *query.Queries, pid int64) (int64, error)
	// Not needed
	Type() string
	// Can be a separate function
	Dialogs() Dialogs
	// Not needed unless I split out field functions
	Fields() Fields
	// This would be IsFieldTypeValid
	IsFieldNameValid(f string) bool
	// Unneeded unless I need riders on fetching content from the db
	Content(q *query.Queries, rid int64) (content, error)
	// Can be done with field-level validators based on type
	IsContentValid(c content) bool
	// Not needed, field status will be stored on the field
	ContentReview(q *query.Queries, rid int64) (contentreview, error)
	// Unneeded, will be done by the RFID
	UpdateField(q *query.Queries, p UpdateFieldParams) error
	// Unneeded, can be broken out into a function or a trait-inheritor
	TitleForSummary(c content) string
	// Unneeded, can be a query
	FieldsForSummary(p FieldsForSummaryParams) ([]FieldForSummary, error)
	// Can be a function or trait-inheritor
	SummaryForQueue(p SummaryForQueueParams) (SummaryForQueue, error)
	// Unneeded, will be a query or just a function
	UpdateFieldStatus(q *query.Queries, p UpdateFieldStatusParams) error
	// Can be a constant by field type
	FieldHelp(e *html.Engine, t, f string) (template.HTML, error)
	// Can be a function or trait inheritor
	RenderFieldForm(e *html.Engine, p RenderFieldFormParams) (template.HTML, error)
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
	Dialogs         Dialogs
	StatusColor     string
	StatusText      string
	Title           string
	Link            string
	AuthorUsername  string
	ReviewerText    template.HTML
	StatusIcon      StatusIcon
	ID              int64
	PID             int64
	ShowPutInReview bool
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
	dialogs := def.Dialogs()
	dialogs.SetPath(p.Request.ID)
	dialogs.PutInReview.Variable = fmt.Sprintf("showReviewDialogForRequest%d", p.Request.ID)
	showPutInReview := CanBePutInReview(
		CanBePutInReviewParams{
			Request:             p.Request,
			ReviewerPermissions: p.ReviewerPermissions,
			PID:                 p.PID,
		},
	)
	// TODO: Make this resilient to a request with an invalid status
	return SummaryForQueue{
		ID:              p.Request.ID,
		PID:             p.Request.PID,
		Title:           title,
		Link:            route.RequestPath(p.Request.ID),
		StatusIcon:      NewStatusIcon(StatusIconParams{Status: p.Request.Status, IconSize: 48, IncludeText: false}),
		StatusColor:     StatusColors[p.Request.Status],
		StatusText:      StatusTexts[p.Request.Status],
		ReviewerText:    reviewerText,
		Dialogs:         dialogs,
		AuthorUsername:  p.PlayerUsername,
		ShowPutInReview: showPutInReview,
	}, nil
}

func (d *DefaultDefinition) FieldHelp(e *html.Engine, t, f string) (template.HTML, error) {
	def, ok := Definitions.Get(t)
	if !ok {
		return template.HTML(""), ErrNoDefinition
	}
	if !IsFieldNameValid(t, f) {
		return template.HTML(""), ErrInvalidInput
	}
	fields := def.Fields()
	return fields.FieldHelp(e, f)
}

func (d *DefaultDefinition) RenderFieldForm(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return template.HTML(""), ErrNoDefinition
	}
	if !IsFieldNameValid(p.Request.Type, p.FieldName) {
		return template.HTML(""), ErrInvalidInput
	}
	fields := def.Fields()
	return fields.RenderForm(e, p)
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

type ForOverviewParams = field.ForOverviewParams

func FieldsForOverview(p field.ForOverviewParams) ([]field.ForOverview, error) {
	fields, ok := NewFieldsByType[p.Request.Type]
	if !ok {
		return []field.ForOverview{}, ErrNoDefinition
	}
	return fields.ForOverview(p), nil
}

// TODO: Make this a standard utility
func ContentBytes(content any) ([]byte, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

var FieldsByType map[string][]string = map[string][]string{
	TypeCharacterApplication: {
		FieldCharacterApplicationName.Name,
		FieldCharacterApplicationGender.Name,
		FieldCharacterApplicationShortDescription.Name,
		FieldCharacterApplicationDescription.Name,
		FieldCharacterApplicationBackstory.Name,
	},
}

var NewFieldsByType map[string]field.Group = map[string]field.Group{
	TypeCharacterApplication: definition.CharacterApplicationFields,
}

type NewParams struct {
	Type string
	PID  int64
}

func NewNew(q *query.Queries, p NewParams) (int64, error) {
	if p.PID == 0 {
		return 0, ErrInvalidInput
	}

	if !IsTypeValid(p.Type) {
		return 0, ErrInvalidType
	}

	fields, ok := NewFieldsByType[p.Type]
	if !ok {
		return 0, ErrInvalidType
	}

	result, err := q.CreateRequest(context.Background(), query.CreateRequestParams{
		PID:  p.PID,
		Type: p.Type,
	})
	if err != nil {
		return 0, err
	}
	rid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, field := range fields.List() {
		if err := q.CreateRequestField(context.Background(), query.CreateRequestFieldParams{
			RID:  rid,
			Type: field.Type,
			// TODO: Store the Default Field Status in its own const
			Status: FieldStatusNotReviewed,
			Value:  "",
		}); err != nil {
			return 0, err
		}
	}

	return rid, nil
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

type NextIncompleteFieldOutput struct {
	Field string
	Last  bool
}

// TODO: Let this error out with no definition
func NextIncompleteField(t string, fieldmap field.Map) (*query.RequestField, bool) {
	fields, ok := NewFieldsByType[t]
	if !ok {
		return nil, false
	}
	return fields.NextIncomplete(fieldmap)
}

type NextUnreviewedFieldOutput struct {
	Field string
	Last  bool
}

func NextUnreviewedField(t string, cr contentreview) (NextUnreviewedFieldOutput, error) {
	definition, ok := Definitions.Get(t)
	if !ok {
		return NextUnreviewedFieldOutput{}, ErrNoDefinition
	}
	fields := definition.Fields()
	return fields.NextUnreviewed(cr)
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

// TODO: Allow an error output here?
func IsFieldTypeValid(t, ft string) bool {
	fields, ok := NewFieldsByType[t]
	if !ok {
		return false
	}
	_, ok = fields.Map()[ft]
	return ok
}

// TODO: Rename to RenderFieldHelp
func FieldHelp(e *html.Engine, t, f string) (template.HTML, error) {
	def, ok := Definitions.Get(t)
	if !ok {
		return template.HTML(""), ErrNoDefinition
	}
	return def.FieldHelp(e, t, f)
}

type RenderFieldDataParams struct {
	Request    *query.Request
	Content    content
	FieldName  string
	FieldValue string
}

func RenderFieldData(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return template.HTML(""), ErrNoDefinition
	}
	if !IsFieldNameValid(p.Request.Type, p.FieldName) {
		return template.HTML(""), ErrInvalidInput
	}
	fields := def.Fields()
	return fields.RenderData(e, p)
}

type RenderFieldFormParams struct {
	Request    *query.Request
	Content    content
	FormID     string
	Path       string
	FieldName  string
	FieldValue string
}

func RenderFieldForm(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return template.HTML(""), ErrNoDefinition
	}
	if !IsFieldNameValid(p.Request.Type, p.FieldName) {
		return template.HTML(""), ErrInvalidInput
	}
	fields := def.Fields()
	return fields.RenderForm(e, p)
}
