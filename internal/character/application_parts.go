package character

import (
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
)

type ApplicationNav struct {
	Link    string
	Label   string
	Current bool
	Ready   bool
}

func MakeApplicationNav(current string, app *queries.CharacterApplicationContent) []ApplicationNav {
	result := []ApplicationNav{}

	result = append(result, ApplicationNav{
		Label:   "Name",
		Link:    routes.RequestFieldPath(app.RID, request.FieldName),
		Current: current == "name",
		Ready:   IsNameValid(app.Name),
	})

	result = append(result, ApplicationNav{
		Label:   "Gender",
		Link:    routes.RequestFieldPath(app.RID, request.FieldGender),
		Current: current == "gender",
		Ready:   IsGenderValid(app.Gender),
	})

	result = append(result, ApplicationNav{
		Label:   "Short Description",
		Link:    routes.RequestFieldPath(app.RID, request.FieldShortDescription),
		Current: current == "sdesc",
		Ready:   IsShortDescriptionValid(app.ShortDescription),
	})

	result = append(result, ApplicationNav{
		Label:   "Description",
		Link:    routes.RequestFieldPath(app.RID, request.FieldDescription),
		Current: current == "description",
		Ready:   IsDescriptionValid(app.Description),
	})

	result = append(result, ApplicationNav{
		Label:   "Backstory",
		Link:    routes.RequestFieldPath(app.RID, request.FieldBackstory),
		Current: current == "backstory",
		Ready:   IsBackstoryValid(app.Backstory),
	})

	return result
}

func IsApplicationReady(app *queries.CharacterApplicationContent) bool {
	if !IsNameValid(app.Name) {
		return false
	}

	if !IsGenderValid(app.Gender) {
		return false
	}

	if !IsShortDescriptionValid(app.ShortDescription) {
		return false
	}

	if !IsDescriptionValid(app.Description) {
		return false
	}

	if !IsBackstoryValid(app.Backstory) {
		return false
	}

	return true
}
