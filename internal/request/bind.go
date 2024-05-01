package request

import (
	"html/template"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/field"
	"petrichormud.com/app/internal/route"
)

type BindFieldViewParams struct {
	Request    *query.Request
	Field      *query.RequestField
	OpenChange *query.OpenRequestChangeRequest
	PID        int64
	Last       bool
}

func NewBindFieldView(e *html.Engine, b fiber.Map, p BindFieldViewParams) (fiber.Map, error) {
	fields, ok := FieldsByType[p.Request.Type]
	if !ok {
		return fiber.Map{}, ErrNoDefinition
	}
	fd, ok := fields.Get(p.Field.Type)
	if !ok {
		return fiber.Map{}, ErrInvalidType
	}
	help, err := fd.RenderHelp(e)
	if err != nil {
		return b, err
	}
	b["Help"] = help

	// TODO: Get this into a utility
	if p.Request.PID == p.PID && p.Request.Status == StatusIncomplete || p.Request.Status == StatusReady {
		form, err := fd.RenderForm(e, p.Field)
		if err != nil {
			return b, err
		}
		b["Form"] = form
	} else {
		data, err := fd.RenderData(e, p.Field)
		if err != nil {
			return b, err
		}
		b["Data"] = data
	}

	b, err = BindDialogs(b, p.Request)
	if err != nil {
		return b, err
	}

	b["FieldLabel"] = fd.Label
	b["FieldDescription"] = fd.Description
	b["RequestFormID"] = FormID

	// TODO: Sort out this being disabled
	b["BackLink"] = route.RequestPath(p.Field.RID)

	b["RequestFormPath"] = route.RequestFieldPath(p.Field.RID, p.Field.Type)
	// TODO: Change this to FieldName
	b["Field"] = p.Field.Type
	b["FieldValue"] = p.Field.Value

	b, err = BindFieldViewActions(e, b, BindFieldViewActionsParams(p))
	if err != nil {
		return fiber.Map{}, err
	}

	b["ChangeRequestConfig"] = BindChangeRequestConfig(BindChangeRequestConfigParams{
		PID:        p.PID,
		OpenChange: p.OpenChange,
		Request:    p.Request,
		Field:      p.Field,
	})
	b["ChangeRequestPath"] = route.RequestChangeRequestFieldPath(p.Request.ID, p.Field.Type)
	if p.OpenChange != nil {
		b["ChangeRequest"] = BindChangeRequest(BindChangeRequestParams{
			PID:    p.PID,
			Change: p.OpenChange,
		})
	}

	return b, nil
}

type BindChangeRequestConfigParams struct {
	OpenChange *query.OpenRequestChangeRequest
	Request    *query.Request
	Field      *query.RequestField
	PID        int64
}

func BindChangeRequestConfig(p BindChangeRequestConfigParams) fiber.Map {
	b := fiber.Map{}
	b["Path"] = route.RequestChangeRequestFieldPath(p.Request.ID, p.Field.Type)
	b["Type"] = p.Field.Type
	if p.OpenChange != nil {
		b["Open"] = BindChangeRequest(BindChangeRequestParams{
			PID:    p.PID,
			Change: p.OpenChange,
		})
	}
	return b
}

type BindFieldViewActionsParams struct {
	Request    *query.Request
	Field      *query.RequestField
	OpenChange *query.OpenRequestChangeRequest
	PID        int64
	Last       bool
}

func BindFieldViewActions(e *html.Engine, b fiber.Map, p BindFieldViewActionsParams) (fiber.Map, error) {
	actions := []template.HTML{}

	if p.Request.Status == StatusInReview && p.Request.RPID == p.PID {
		if p.OpenChange == nil {
			change, err := partial.Render(e, partial.RenderParams{
				Template: partial.RequestFieldActionChangeRequest,
			})
			if err != nil {
				return b, err
			}
			actions = append(actions, change)
		}

		reject, err := partial.Render(e, partial.RenderParams{
			Template: partial.RequestFieldActionReject,
		})
		if err != nil {
			return b, err
		}
		actions = append(actions, reject)

		text := "Approve"
		if p.OpenChange != nil {
			if p.Last {
				text = "Finish"
			} else {
				text = "Next"
			}
		} else {
			text = "Approve"
		}
		review, err := partial.Render(e, partial.RenderParams{
			Template: partial.RequestFieldActionReview,
			Bind: fiber.Map{
				"Path": route.RequestFieldStatusPath(p.Request.ID, p.Field.Type),
				"Text": text,
			},
		})
		if err != nil {
			return b, err
		}
		actions = append(actions, review)
	}

	// TODO: Bind this to the same function that determines if we show the form or not
	if p.Request.PID == p.PID && p.Request.Status == StatusIncomplete || p.Request.Status == StatusReady {
		text := "Next"
		if p.Request.Status == StatusReady {
			text = "Update"
		}
		if p.Request.Status == StatusIncomplete && p.Last {
			text = "Finish"
		}
		// TODO: Set this up so the button is disabled if the field is incomplete
		update, err := partial.Render(e, partial.RenderParams{
			Template: partial.RequestFieldActionUpdate,
			Bind: fiber.Map{
				"Form": FormID,
				"Text": text,
			},
		})
		if err != nil {
			return b, err
		}
		actions = append(actions, update)
	}

	b["Actions"] = actions
	return b, nil
}

type BindOverviewParams struct {
	Request  *query.Request
	FieldMap field.Map
	PID      int64
}

func BindOverview(e *html.Engine, b fiber.Map, p BindOverviewParams) (fiber.Map, error) {
	b["PageHeader"] = fiber.Map{
		"Title": Title(p.Request.Type, p.FieldMap),
	}
	// TODO: Build a utility for this
	b["Status"] = fiber.Map{
		"StatusIcon": NewStatusIcon(StatusIconParams{Status: p.Request.Status, IconSize: 48, IncludeText: true, TextSize: "text-xl"}),
	}

	b, err := BindOverviewActions(e, b, BindOverviewActionsParams{
		PID:     p.PID,
		Request: p.Request,
	})
	if err != nil {
		return b, err
	}

	return b, nil
}

type BindOverviewActionsParams struct {
	Request  *query.Request
	FieldMap field.Map
	PID      int64
}

func BindOverviewActions(e *html.Engine, b fiber.Map, p BindOverviewActionsParams) (fiber.Map, error) {
	actions := []template.HTML{}

	if p.Request.PID == p.PID {
		cancel, err := partial.Render(e, partial.RenderParams{
			Template: partial.RequestOverviewActionCancel,
		})
		if err != nil {
			return b, err
		}
		actions = append(actions, cancel)

		if p.Request.Status == StatusReady {
			submit, err := partial.Render(e, partial.RenderParams{
				Template: partial.RequestOverviewActionSubmit,
			})
			if err != nil {
				return b, err
			}
			actions = append(actions, submit)
		}
	}

	if p.Request.Status == StatusInReview && p.Request.RPID == p.PID {
		reject, err := partial.Render(e, partial.RenderParams{
			Template: partial.RequestOverviewActionReject,
		})
		if err != nil {
			return b, err
		}
		actions = append(actions, reject)

		allApproved := true
		for _, field := range p.FieldMap {
			if field.Status != FieldStatusApproved {
				allApproved = false
			}
		}

		if allApproved {
			approve, err := partial.Render(e, partial.RenderParams{
				Template: partial.RequestOverviewActionApprove,
			})
			if err != nil {
				return b, err
			}
			actions = append(actions, approve)
		} else {
			review, err := partial.Render(e, partial.RenderParams{
				Template: partial.RequestOverviewActionReview,
			})
			if err != nil {
				return b, err
			}
			actions = append(actions, review)
		}
	}

	b["Actions"] = actions
	return b, nil
}

type BindChangeRequestParams struct {
	Change *query.OpenRequestChangeRequest
	PID    int64
}

func BindChangeRequest(p BindChangeRequestParams) fiber.Map {
	b := fiber.Map{
		"Text": p.Change.Text,
		"Path": route.RequestChangeRequestPath(p.Change.ID),
	}

	if p.Change.PID == p.PID {
		b["ShowDeleteAction"] = true
		b["ShowEditAction"] = true
	}

	return b
}
