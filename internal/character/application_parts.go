package character

import (
	"strconv"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

type ApplicationPart struct {
	Link    string
	Label   string
	Current bool
	Ready   bool
}

func MakeApplicationParts(current string, app *queries.CharacterApplicationContent) []ApplicationPart {
	result := []ApplicationPart{}

	result = append(result, ApplicationPart{
		Label:   "Name",
		Link:    routes.CharacterApplicationNamePath(strconv.FormatInt(app.RID, 10)),
		Current: current == "name",
		Ready:   IsNameValid(app.Name),
	})

	result = append(result, ApplicationPart{
		Label:   "Gender",
		Link:    routes.CharacterApplicationGenderPath(strconv.FormatInt(app.RID, 10)),
		Current: current == "gender",
		Ready:   IsGenderValid(app.Gender),
	})

	result = append(result, ApplicationPart{
		Label:   "Short Description",
		Link:    routes.CharacterApplicationShortDescriptionPath(strconv.FormatInt(app.RID, 10)),
		Current: current == "sdesc",
		Ready:   IsShortDescriptionValid(app.ShortDescription),
	})

	result = append(result, ApplicationPart{
		Label:   "Description",
		Link:    routes.CharacterApplicationDescriptionPath(strconv.FormatInt(app.RID, 10)),
		Current: current == "description",
		Ready:   IsDescriptionValid(app.Description),
	})

	result = append(result, ApplicationPart{
		Label:   "Backstory",
		Link:    routes.CharacterApplicationBackstoryPath(strconv.FormatInt(app.RID, 10)),
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
