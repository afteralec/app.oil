package request

import (
	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/actor"
	"petrichormud.com/app/internal/bind"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

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

// TODO: Add ViewedByAdmin and maybe have ViewedByPlayer be ViewedByOwner
func BindViewedBy(b fiber.Map, p BindViewedByParams) fiber.Map {
	b["ViewedByPlayer"] = p.Request.PID == p.PID
	b["ViewedByReviewer"] = p.Request.RPID == p.PID

	return b
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
