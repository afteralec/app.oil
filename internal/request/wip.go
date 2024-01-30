package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"petrichormud.com/app/internal/queries"
)

const errInvalidType string = "invalid type"

var ErrInvalidType error = errors.New(errInvalidType)

func Content(q *queries.Queries, req *queries.Request) (map[string]string, error) {
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
	Request *queries.Request
	Field   string
	Value   string
	PID     int64
}

func UpdateField(q *queries.Queries, p UpdateFieldParams) error {
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

type GetSummaryFieldsParams struct {
	Request *queries.Request
	Content map[string]string
	PID     int64
}

// TODO: Get this built into the Definition
func GetSummaryFields(p GetSummaryFieldsParams) []SummaryField {
	if p.Request.Type == TypeCharacterApplication {
		var basePathSB strings.Builder
		fmt.Fprintf(&basePathSB, "/requests/%d", p.Request.ID)
		basePath := basePathSB.String()

		var namePathSB strings.Builder
		fmt.Fprintf(&namePathSB, "%s/%s", basePath, FieldName)

		var genderPathSB strings.Builder
		fmt.Fprintf(&genderPathSB, "%s/%s", basePath, FieldGender)

		var shortDescriptionPathSB strings.Builder
		fmt.Fprintf(&shortDescriptionPathSB, "%s/%s", basePath, FieldShortDescription)

		var descriptionPathSB strings.Builder
		fmt.Fprintf(&descriptionPathSB, "%s/%s", basePath, FieldDescription)

		var backstoryPathSB strings.Builder
		fmt.Fprintf(&backstoryPathSB, "%s/%s", basePath, FieldBackstory)

		allowEdit := p.Request.PID == p.PID

		return []SummaryField{
			{
				Label:     "Name",
				Content:   p.Content[FieldName],
				AllowEdit: allowEdit,
				Path:      namePathSB.String(),
			},
			{
				Label:     "Gender",
				Content:   p.Content[FieldGender],
				AllowEdit: allowEdit,
				Path:      genderPathSB.String(),
			},
			{
				Label:     "Short Description",
				Content:   p.Content[FieldShortDescription],
				AllowEdit: allowEdit,
				Path:      shortDescriptionPathSB.String(),
			},
			{
				Label:     "Description",
				Content:   p.Content[FieldDescription],
				AllowEdit: allowEdit,
				Path:      descriptionPathSB.String(),
			},
			{
				Label:     "Backstory",
				Content:   p.Content[FieldBackstory],
				AllowEdit: allowEdit,
				Path:      backstoryPathSB.String(),
			},
		}
	}

	return []SummaryField{}
}
