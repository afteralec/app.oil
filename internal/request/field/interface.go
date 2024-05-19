package field

import (
	"html/template"
	"maps"
	"slices"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/change"
	"petrichormud.com/app/internal/request/status"
	"petrichormud.com/app/internal/route"
	"petrichormud.com/app/internal/validate"
)

type Renderer interface {
	Render(e *html.Engine, field *query.RequestField, template string) (template.HTML, error)
}

type DefaultRenderer struct{}

func (r *DefaultRenderer) Render(e *html.Engine, field *query.RequestField, template string) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: template,
		Bind: fiber.Map{
			"FieldValue": field.Value,
			"FormID":     "request-form",
			"Path":       route.RequestFieldPath(field.RID, field.Type),
		},
	})
}

// TODO: Rename this to Definition or Configuration - I prefer Definition so far
type Field struct {
	Validator      validate.StringValidator
	FormRenderer   Renderer
	Type           string
	For            string
	Label          string
	Description    string
	Help           string
	Data           string
	Form           string
	SubfieldConfig SubfieldConfig
}

type SubfieldConfig struct {
	MinValues int
	MaxValues int
	Require   bool
}

func NewSubfieldConfig(min, max int) SubfieldConfig {
	return SubfieldConfig{
		Require:   true,
		MinValues: min,
		MaxValues: max,
	}
}

func (f *Field) IsValid(v string) bool {
	return f.Validator.IsValid(v)
}

func (f *Field) RenderHelp(e *html.Engine) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: f.Help,
	})
}

func (f *Field) RenderData(e *html.Engine, field *query.RequestField) (template.HTML, error) {
	return partial.Render(e, partial.RenderParams{
		Template: f.Data,
		Bind: fiber.Map{
			"FieldValue": field.Value,
		},
	})
}

func (f *Field) RenderForm(e *html.Engine, field *query.RequestField) (template.HTML, error) {
	return f.FormRenderer.Render(e, field, f.Form)
}

func (f *Field) ForPlayer() bool {
	return f.For == ForPlayer
}

func (f *Field) ForReviewer() bool {
	return f.For == ForReviewer
}

type ForOverview struct {
	// TODO: Get this into a discrete type instead of a fiber Map?
	ChangeRequestConfig     fiber.Map
	Help                    template.HTML
	Type                    string
	Label                   string
	Value                   string
	Path                    string
	AllowEdit               bool
	IsApproved              bool
	ShowRequestChangeAction bool
}

type ForOverviewParams struct {
	Request       *query.Request
	FieldMap      Map
	ChangeMap     map[int64]query.RequestChangeRequest
	OpenChangeMap map[int64]query.OpenRequestChangeRequest
	PID           int64
}

func (f *Field) ForOverview(e *html.Engine, p ForOverviewParams) ForOverview {
	v := ""
	field, ok := p.FieldMap[f.Type]
	if ok {
		v = field.Value
	}

	// TODO: Build a utility for this
	allowEdit := p.Request.PID == p.PID
	if p.Request.Status != status.Incomplete && p.Request.Status != status.Ready && p.Request.Status != status.Reviewed {
		allowEdit = false
	}

	help, err := f.RenderHelp(e)
	if err != nil {
		// TODO: Handle this error
		help = template.HTML("")
	}

	overview := ForOverview{
		Help:                    help,
		Type:                    f.Type,
		Label:                   f.Label,
		Value:                   v,
		Path:                    route.RequestFieldPath(p.Request.ID, f.Type),
		AllowEdit:               allowEdit,
		IsApproved:              field.Status == StatusApproved,
		ShowRequestChangeAction: p.PID == p.Request.RPID && p.Request.Status == status.InReview,
	}

	bcp := change.BindConfigParams{
		Request: p.Request,
		Field:   &field,
		PID:     p.PID,
	}

	openchange, ok := p.OpenChangeMap[field.ID]
	if ok {
		bcp.OpenChange = &openchange
	}
	ch, ok := p.ChangeMap[field.ID]
	if ok {
		bcp.Change = &ch
	}
	overview.ChangeRequestConfig = change.BindConfig(bcp)

	return overview
}

type fieldBuilder struct {
	Field Field
}

func FieldBuilder() *fieldBuilder {
	return new(fieldBuilder)
}

func (b *fieldBuilder) Type(t string) *fieldBuilder {
	b.Field.Type = t
	return b
}

func (b *fieldBuilder) For(f string) *fieldBuilder {
	b.Field.For = f
	return b
}

func (b *fieldBuilder) Label(label string) *fieldBuilder {
	b.Field.Label = label
	return b
}

func (b *fieldBuilder) Description(description string) *fieldBuilder {
	b.Field.Description = description
	return b
}

func (b *fieldBuilder) Help(help string) *fieldBuilder {
	b.Field.Help = help
	return b
}

func (b *fieldBuilder) Data(data string) *fieldBuilder {
	b.Field.Data = data
	return b
}

func (b *fieldBuilder) Form(form string) *fieldBuilder {
	b.Field.Form = form
	return b
}

func (b *fieldBuilder) FormRenderer(r Renderer) *fieldBuilder {
	b.Field.FormRenderer = r
	return b
}

func (b *fieldBuilder) Validator(validator validate.StringValidator) *fieldBuilder {
	b.Field.Validator = validator
	return b
}

func (b *fieldBuilder) SubfieldConfig(config SubfieldConfig) *fieldBuilder {
	b.Field.SubfieldConfig = config
	return b
}

func (b *fieldBuilder) Build() Field {
	// TODO: Allow fields to have default values
	// TODO: Validate that the field is being built with all of its needed parts
	return b.Field
}

type Group struct {
	fields map[string]Field
	list   []Field
}

func (fg *Group) Map() map[string]Field {
	return maps.Clone(fg.fields)
}

func (fg *Group) List() []Field {
	return slices.Clone(fg.list)
}

func (fg *Group) Get(ft string) (Field, bool) {
	field, ok := fg.fields[ft]
	return field, ok
}

type NextIncompleteOutput struct {
	Field *query.RequestField
	Last  bool
}

func (f *Group) NextIncomplete(fields Map) (NextIncompleteOutput, error) {
	for i, fd := range f.list {
		field, ok := fields[fd.Type]
		if !ok {
			continue
		}
		if len(field.Value) == 0 {
			return NextIncompleteOutput{
				Field: &field,
				Last:  i == len(f.list)-1,
			}, nil
		}
	}
	return NextIncompleteOutput{}, nil
}

type NextUnreviewedOutput struct {
	Field *query.RequestField
	Last  bool
}

func (f *Group) NextUnreviewed(fields Map) (NextUnreviewedOutput, error) {
	for i, fd := range f.list {
		field, ok := fields[fd.Type]
		if !ok {
			continue
		}
		if field.Status == StatusNotReviewed {
			return NextUnreviewedOutput{
				Field: &field,
				Last:  i == len(f.list)-1,
			}, nil
		}
	}
	return NextUnreviewedOutput{}, nil
}

func (f *Group) ForOverview(e *html.Engine, p ForOverviewParams) []ForOverview {
	result := []ForOverview{}
	for _, field := range f.list {
		result = append(result, field.ForOverview(e, p))
	}
	return result
}

func NewGroup(fields []Field) Group {
	m := map[string]Field{}
	l := []Field{}
	for _, field := range fields {
		m[field.Type] = field
		l = append(l, field)
	}
	return Group{
		fields: m,
		list:   l,
	}
}
