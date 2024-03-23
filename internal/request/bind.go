package request

import (
	"html/template"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

type BindFieldViewParams struct {
	Request               *query.Request
	Content               content
	FieldName             string
	CurrentChangeRequests []query.RequestChangeRequest
	PID                   int64
	Last                  bool
}

func BindFieldView(e *html.Engine, b fiber.Map, p BindFieldViewParams) (fiber.Map, error) {
	help, err := FieldHelp(e, p.Request.Type, p.FieldName)
	if err != nil {
		return b, err
	}
	b["Help"] = help

	// TODO: Get this into a utility
	if p.Request.PID == p.PID && p.Request.Status == StatusIncomplete || p.Request.Status == StatusReady {
		fieldValue, ok := p.Content.Value(p.FieldName)
		if !ok {
			fieldValue = ""
		}

		form, err := RenderFieldForm(e, RenderFieldFormParams{
			Request:    p.Request,
			Content:    p.Content,
			FieldName:  p.FieldName,
			FieldValue: fieldValue,
			FormID:     FormID,
		})
		if err != nil {
			return b, err
		}
		b["Form"] = form
	} else {
		fieldValue, ok := p.Content.Value(p.FieldName)
		if !ok {
			fieldValue = ""
		}

		data, err := RenderFieldData(e, RenderFieldDataParams{
			Request:    p.Request,
			Content:    p.Content,
			FieldName:  p.FieldName,
			FieldValue: fieldValue,
		})
		if err != nil {
			return b, err
		}
		b["Data"] = data
	}

	b, err = BindDialogs(b, BindDialogsParams{
		Request: p.Request,
	})
	if err != nil {
		return b, err
	}

	label, description := GetFieldLabelAndDescription(p.Request.Type, p.FieldName)
	b["FieldLabel"] = label
	b["FieldDescription"] = description

	b["RequestFormID"] = FormID

	// TODO: Sort out this being disabled
	b["BackLink"] = route.RequestPath(p.Request.ID)

	b["RequestFormPath"] = route.RequestFieldPath(p.Request.ID, p.FieldName)
	// TODO: Change this to FieldName
	b["Field"] = p.FieldName

	// TODO: Consolidate this with the above
	fieldValue, ok := p.Content.Value(p.FieldName)
	if ok {
		b["FieldValue"] = fieldValue
	} else {
		b["FieldValue"] = ""
	}

	BindFieldViewActions(e, b, BindFieldViewActionsParams{
		PID:                   p.PID,
		Request:               p.Request,
		CurrentChangeRequests: p.CurrentChangeRequests,
		FieldName:             p.FieldName,
		Last:                  p.Last,
	})

	// TODO: Move this to a utility
	b["ChangeRequestPath"] = route.RequestChangeRequestFieldPath(p.Request.ID, p.FieldName)
	if len(p.CurrentChangeRequests) == 1 {
		b["ChangeRequest"] = BindChangeRequest(BindChangeRequestParams{
			PID:           p.PID,
			ChangeRequest: &p.CurrentChangeRequests[0],
		})
	}

	return b, nil
}

type BindFieldViewActionsParams struct {
	Request               *query.Request
	FieldName             string
	CurrentChangeRequests []query.RequestChangeRequest
	PID                   int64
	Last                  bool
}

func BindFieldViewActions(e *html.Engine, b fiber.Map, p BindFieldViewActionsParams) (fiber.Map, error) {
	actions := []template.HTML{}

	if p.Request.Status == StatusInReview && p.Request.RPID == p.PID {
		// TODO: Put this in a utility
		if len(p.CurrentChangeRequests) == 0 {
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
		if len(p.CurrentChangeRequests) > 0 {
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
				"Path": route.RequestFieldStatusPath(p.Request.ID, p.FieldName),
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

type BindChangeRequestParams struct {
	ChangeRequest *query.RequestChangeRequest
	PID           int64
}

func BindChangeRequest(p BindChangeRequestParams) fiber.Map {
	b := fiber.Map{
		"Text": p.ChangeRequest.Text,
		"Path": route.RequestChangeRequestPath(p.ChangeRequest.ID),
	}

	if p.ChangeRequest.PID == p.PID && !p.ChangeRequest.Locked && !p.ChangeRequest.Old {
		b["ShowDeleteAction"] = true
		b["ShowEditAction"] = true
	}

	return b
}
