package request

import (
	"errors"
	"html/template"

	html "github.com/gofiber/template/html/v2"
	"petrichormud.com/app/internal/query"
)

var ErrInvalidInput error = errors.New("field value didn't pass validation")

type Fields struct {
	Map  map[string]Field
	List []Field
}

func NewFields(fields []Field) Fields {
	fieldsMap := map[string]Field{}
	for _, field := range fields {
		fieldsMap[field.Name] = field
	}
	return Fields{
		List: fields,
		Map:  fieldsMap,
	}
}

func (f *Fields) Update(q *query.Queries, p UpdateFieldParams) error {
	field, ok := f.Map[p.FieldName]
	if !ok {
		return ErrInvalidInput
	}
	return field.Update(q, p)
}

func (f *Fields) UpdateStatus(q *query.Queries, p UpdateFieldStatusParams) error {
	field, ok := f.Map[p.FieldName]
	if !ok {
		return ErrInvalidInput
	}
	return field.UpdateStatus(q, p)
}

func (f *Fields) NextIncomplete(c content) (string, bool) {
	for i, field := range f.List {
		value, ok := c.Value(field.Name)
		if !ok {
			continue
		}
		if len(value) == 0 {
			return field.Name, i == len(f.List)-1
		}
	}
	return "", false
}

func (f *Fields) NextUnreviewed(cr contentreview) (NextUnreviewedFieldOutput, error) {
	for i, field := range f.List {
		status, ok := cr.Status(field.Name)
		if !ok {
			continue
		}
		if status == FieldStatusNotReviewed {
			return NextUnreviewedFieldOutput{
				Field: field.Name,
				Last:  i == len(f.List)-1,
			}, nil
		}
	}

	return NextUnreviewedFieldOutput{
		Field: "",
		Last:  false,
	}, nil
}

func (f *Fields) IsFieldNameValid(name string) bool {
	_, ok := f.Map[name]
	return ok
}

func (f *Fields) IsFieldValueValid(name, value string) bool {
	field, ok := f.Map[name]
	if !ok {
		return false
	}
	return field.IsValueValid(value)
}

func (f *Fields) ForSummary(p FieldsForSummaryParams) []FieldForSummary {
	result := []FieldForSummary{}
	for _, field := range f.List {
		result = append(result, field.ForSummary(p))
	}
	return result
}

func (f *Fields) FieldHelp(e *html.Engine, name string) (template.HTML, error) {
	field, ok := f.Map[name]
	if !ok {
		// TODO: ErrInvalidField
		return template.HTML(""), ErrInvalidInput
	}
	return field.RenderHelp(e)
}

func (f *Fields) RenderData(e *html.Engine, p RenderFieldDataParams) (template.HTML, error) {
	field, ok := f.Map[p.FieldName]
	if !ok {
		// TODO: ErrInvalidField
		return template.HTML(""), ErrInvalidInput
	}
	return field.RenderData(e, p)
}

func (f *Fields) RenderForm(e *html.Engine, p RenderFieldFormParams) (template.HTML, error) {
	field, ok := f.Map[p.FieldName]
	if !ok {
		// TODO: ErrInvalidField
		return template.HTML(""), ErrInvalidInput
	}
	return field.RenderForm(e, p)
}
