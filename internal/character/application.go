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
	if !req.New {
		return routes.CharacterApplicationPath(strconv.FormatInt(req.ID, 10))
	}

	if len(app.Name) == 0 {
		return routes.CharacterApplicationNamePath(strid)
	}

	if len(app.Gender) == 0 {
		return routes.CharacterApplicationGenderPath(strid)
	}

	// TODO: Rename this to ShortDescription
	if len(app.Sdesc) == 0 {
		return routes.CharacterApplicationShortDescriptionPath(strid)
	}

	if len(app.Description) == 0 {
		return routes.CharacterApplicationDescriptionPath(strid)
	}

	if len(app.Backstory) == 0 {
		return routes.CharacterApplicationBackstoryPath(strid)
	}

	return routes.CharacterApplicationPath(strconv.FormatInt(req.ID, 10))
}
