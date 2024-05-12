package change

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

type BindParams struct {
	Field      *query.RequestField
	OpenChange *query.OpenRequestChangeRequest
	Change     *query.RequestChangeRequest
	PID        int64
}

func Bind(p BindParams) fiber.Map {
	value := ""
	text := ""
	if p.OpenChange != nil {
		if p.OpenChange.Value != p.Field.Value {
			value = p.OpenChange.Value
		}
		text = p.OpenChange.Text
	} else if p.Change != nil {
		if p.Change.Value != p.Field.Value {
			value = p.Change.Value
		}
		text = p.Change.Text
	}
	var id int64 = 0
	if p.OpenChange != nil {
		id = p.OpenChange.ID
	} else if p.Change != nil {
		id = p.Change.ID
	}

	b := fiber.Map{
		"Text": text,
		"Path": route.RequestChangeRequestPath(id),
	}

	if len(value) > 0 {
		b["FieldValue"] = value
	}

	if p.OpenChange != nil && p.OpenChange.PID == p.PID {
		b["ShowDeleteAction"] = true
		b["ShowEditAction"] = true
	}

	return b
}

type BindConfigParams struct {
	OpenChange *query.OpenRequestChangeRequest
	Change     *query.RequestChangeRequest
	Request    *query.Request
	Field      *query.RequestField
	PID        int64
}

func BindConfig(p BindConfigParams) fiber.Map {
	b := fiber.Map{}
	b["Path"] = route.RequestChangeRequestFieldPath(p.Request.ID, p.Field.Type)
	b["Type"] = p.Field.Type
	if p.OpenChange != nil {
		b["Open"] = Bind(BindParams{
			PID:        p.PID,
			OpenChange: p.OpenChange,
			Field:      p.Field,
		})
	}
	if p.Change != nil {
		b["Change"] = Bind(BindParams{
			PID:    p.PID,
			Change: p.Change,
			Field:  p.Field,
		})
	}
	return b
}
