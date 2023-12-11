package character

import (
	"strconv"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

type ApplicationSummary struct {
	Link string
	Name string
	ID   int64
}

const DefaultApplicationSummaryName = "Unnamed"

func NewSummaryFromApplication(req *queries.Request, app *queries.CharacterApplicationContent) ApplicationSummary {
	name := app.Name
	if len(app.Name) == 0 {
		name = DefaultApplicationSummaryName
	}
	return ApplicationSummary{
		Link: GetApplicationFlowLink(req, app),
		ID:   req.ID,
		Name: name,
	}
}

func GetApplicationFlowLink(req *queries.Request, app *queries.CharacterApplicationContent) string {
	strid := strconv.FormatInt(req.ID, 10)

	if !IsNameValid(app.Name) {
		return routes.CharacterApplicationNamePath(strid)
	}

	if !IsGenderValid(app.Gender) {
		return routes.CharacterApplicationGenderPath(strid)
	}

	if !IsShortDescriptionValid(app.ShortDescription) {
		return routes.CharacterApplicationShortDescriptionPath(strid)
	}

	if !IsDescriptionValid(app.Description) {
		return routes.CharacterApplicationDescriptionPath(strid)
	}

	if !IsBackstoryValid(app.Backstory) {
		return routes.CharacterApplicationBackstoryPath(strid)
	}

	return routes.CharacterApplicationReviewPath(strconv.FormatInt(req.ID, 10))
}
