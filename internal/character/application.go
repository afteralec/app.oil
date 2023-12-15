package character

import (
	"fmt"
	"strconv"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/request"
	"petrichormud.com/app/internal/routes"
)

type ReviewDialogData struct {
	Path     string
	Variable string
}

type ApplicationSummary struct {
	ReviewDialog     ReviewDialogData
	Status           string
	Link             string
	Name             string
	Author           string
	Reviewer         string
	ID               int64
	RPID             int64
	StatusIncomplete bool
	StatusReady      bool
	StatusSubmitted  bool
	StatusInReview   bool
	StatusApproved   bool
	StatusReviewed   bool
	StatusRejected   bool
	StatusArchived   bool
	StatusCanceled   bool
	Reviewed         bool
}

const DefaultApplicationSummaryName = "Unnamed"

func NewSummaryFromApplication(p *queries.Player, reviewer string, req *queries.Request, app *queries.CharacterApplicationContent) ApplicationSummary {
	name := app.Name
	if len(app.Name) == 0 {
		name = DefaultApplicationSummaryName
	}

	reviewed := len(reviewer) > 0

	return ApplicationSummary{
		Reviewed: reviewed,
		ReviewDialog: ReviewDialogData{
			Path:     "/put/character/application/in/review/test/path",
			Variable: fmt.Sprintf("showReviewDialogFor%s%s", app.Name, p.Username),
		},
		Status:           req.Status,
		StatusIncomplete: req.Status == request.StatusIncomplete,
		StatusReady:      req.Status == request.StatusReady,
		StatusSubmitted:  req.Status == request.StatusSubmitted,
		StatusInReview:   req.Status == request.StatusInReview,
		StatusApproved:   req.Status == request.StatusApproved,
		StatusReviewed:   req.Status == request.StatusReviewed,
		StatusRejected:   req.Status == request.StatusRejected,
		StatusArchived:   req.Status == request.StatusArchived,
		StatusCanceled:   req.Status == request.StatusCanceled,
		Link:             GetApplicationLink(req, app),
		ID:               req.ID,
		Name:             name,
		Author:           p.Username,
		Reviewer:         reviewer,
	}
}

func GetApplicationLink(req *queries.Request, app *queries.CharacterApplicationContent) string {
	if req.Status == request.StatusSubmitted {
		return routes.CharacterApplicationSubmittedPath(strconv.FormatInt(req.ID, 10))
	}

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
