package request

import (
	"context"
	"encoding/json"
	"errors"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/queries"
)

// TODO: Find a way to get each content extractor into the definition of the request type
func GetContent(qtx *queries.Queries, req *queries.Request) (map[string]string, error) {
	var b []byte
	m := map[string]string{}

	switch req.Type {
	case TypeCharacterApplication:
		app, err := qtx.GetCharacterApplicationContentForRequest(context.Background(), req.ID)
		if err != nil {
			return m, err
		}

		b, err = json.Marshal(app)
		if err != nil {
			return m, err
		}
	default:
		return m, errors.New("invalid type")
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]string{}, err
	}

	return m, nil
}

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

type ResolveViewOutput struct {
	Bind   fiber.Map
	View   string
	Layout string
}

type ResolveViewParams struct {
	Bind    fiber.Map
	Request *queries.Request
	PID     int64
}

func ResolveView(qtx *queries.Queries, p ResolveViewParams) (ResolveViewOutput, error) {
	content, err := GetContent(qtx, p.Request)
	if err != nil {
		return ResolveViewOutput{}, err
	}

	switch p.Request.Status {
	case StatusIncomplete:
		field, last := GetNextIncompleteField(p.Request.Type, content)
		view := GetView(p.Request.Type, field)

		label, description := GetFieldLabelAndDescription(p.Request.Type, field)
		p.Bind["Header"] = label
		p.Bind["SubHeader"] = description

		// TODO: Move to this
		p.Bind["Label"] = label
		p.Bind["Description"] = description

		// TODO: Put this in a constant
		p.Bind["RequestFormID"] = FormID

		if last {
			p.Bind["UpdateButtonText"] = "Finish"
		} else {
			p.Bind["UpdateButtonText"] = "Next"
		}

		return ResolveViewOutput{
			Bind:   p.Bind,
			View:   view,
			Layout: "",
		}, nil
	default:
		return ResolveViewOutput{}, errors.New("invalid status")
	}
}
