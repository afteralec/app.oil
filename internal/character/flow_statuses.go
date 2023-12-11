package character

import (
	"strconv"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

type ApplicationPartStatus struct {
	Link    string
	Label   string
	Current bool
	Ready   bool
}

func MakeApplicationPartStatuses(current string, app *queries.CharacterApplicationContent) []ApplicationPartStatus {
	result := []ApplicationPartStatus{}

	result = append(result, ApplicationPartStatus{
		Label:   "Name",
		Link:    routes.CharacterApplicationNamePath(strconv.FormatInt(app.Rid, 10)),
		Current: current == "name",
		Ready:   IsNameValid(app.Name),
	})

	result = append(result, ApplicationPartStatus{
		Label:   "Gender",
		Link:    routes.CharacterApplicationGenderPath(strconv.FormatInt(app.Rid, 10)),
		Current: current == "gender",
		Ready:   IsGenderValid(app.Gender),
	})

	result = append(result, ApplicationPartStatus{
		Label:   "Short Description",
		Link:    routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(app.Rid, 10)),
		Current: current == "sdesc",
		Ready:   IsShortDescriptionValid(app.ShortDescription),
	})

	result = append(result, ApplicationPartStatus{
		Label:   "Description",
		Link:    routes.CharacterApplicationDescriptionPath(strconv.FormatInt(app.Rid, 10)),
		Current: current == "description",
		Ready:   IsDescriptionValid(app.Description),
	})

	result = append(result, ApplicationPartStatus{
		Label:   "Backstory",
		Link:    routes.CharacterApplicationBackstoryPath(strconv.FormatInt(app.Rid, 10)),
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
