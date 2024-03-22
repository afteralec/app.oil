package request

import (
	"html/template"

	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/partial"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

type BindFieldViewParams struct {
	Request       *query.Request
	ChangeRequest *query.RequestChangeRequest
	Content       content
	FieldName     string
	PID           int64
	Last          bool
}

func BindFieldView(e *html.Engine, b fiber.Map, p BindFieldViewParams) (fiber.Map, error) {
	if p.Request.PID == p.PID {
		b["ShowCancelAction"] = true

		switch p.Request.Status {
		case StatusIncomplete:
			b["AllowEdit"] = true
		case StatusReady:
			b["ShowSubmitAction"] = true
			b["AllowEdit"] = true
		}
	}

	b, err := BindDialogs(b, BindDialogsParams{
		Request: p.Request,
	})
	if err != nil {
		return b, err
	}

	label, description := GetFieldLabelAndDescription(p.Request.Type, p.FieldName)
	b["FieldLabel"] = label
	b["FieldDescription"] = description
	if p.Last {
		b["UpdateButtonText"] = "Finish"
	} else {
		b["UpdateButtonText"] = "Next"
	}

	b["RequestFormID"] = FormID

	b["UpdateButtonText"] = "Update"
	b["BackLink"] = route.RequestPath(p.Request.ID)

	b["RequestFormPath"] = route.RequestFieldPath(p.Request.ID, p.FieldName)
	// TODO: Change this to FieldName
	b["Field"] = p.FieldName

	fieldValue, ok := p.Content.Value(p.FieldName)
	if ok {
		b["FieldValue"] = fieldValue
	} else {
		b["FieldValue"] = ""
	}

	b = BindGenderRadioGroup(b, BindGenderRadioGroupParams{
		Content: p.Content,
		Name:    "value",
	})

	actions := []template.HTML{}

	if p.Request.Status == StatusInReview && p.Request.RPID == p.PID {
		// TODO: Put this in a utility
		if p.ChangeRequest == nil {
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
		if p.ChangeRequest == nil {
			text = "Approve"
		} else if p.Last {
			text = "Finish"
		} else {
			text = "Next"
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

	b["Actions"] = actions
	b["ChangeRequestPath"] = route.RequestChangeRequestFieldPath(p.Request.ID, p.FieldName)

	if p.ChangeRequest != nil {
		b["ChangeRequest"] = BindChangeRequest(BindChangeRequestParams{
			PID:           p.PID,
			ChangeRequest: p.ChangeRequest,
		})
	}

	return b, nil
}

type BindGenderRadioGroupParams struct {
	Content content
	Name    string
}

// TODO: Put this behind a Character Applications, Characters or actor package instead?
func BindGenderRadioGroup(b fiber.Map, p BindGenderRadioGroupParams) fiber.Map {
	gender, ok := p.Content.Value(FieldCharacterApplicationGender.Name)
	if !ok {
		return fiber.Map{}
	}
	b["GenderRadioGroup"] = []bind.Radio{
		{
			ID:       "edit-request-character-application-gender-non-binary",
			Name:     p.Name,
			Variable: "gender",
			Value:    actor.GenderNonBinary,
			Label:    "Non-Binary",
			Active:   gender == actor.GenderNonBinary,
		},
		{
			ID:       "edit-request-character-application-gender-female",
			Name:     p.Name,
			Variable: "gender",
			Value:    actor.GenderFemale,
			Label:    "Female",
			Active:   gender == actor.GenderFemale,
		},
		{
			ID:       "edit-request-character-application-gender-male",
			Name:     p.Name,
			Variable: "gender",
			Value:    actor.GenderMale,
			Label:    "Male",
			Active:   gender == actor.GenderMale,
		},
	}
	return b
}

type BindViewedByParams struct {
	Request *query.Request
	PID     int64
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
