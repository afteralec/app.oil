package request

import (
	"fmt"
	"strconv"

	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
)

// TODO: This is a pure dump from the character module; this badly needs cleaning up

type ReviewDialogData struct {
	Path     string
	Variable string
}

type ApplicationSummary struct {
	StatusIcon       StatusIcon
	ReviewDialog     ReviewDialogData
	Status           string
	StatusColor      string
	StatusText       string
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
			Path:     routes.PutCharacterApplicationInReviewPath(strconv.FormatInt(req.ID, 10)),
			Variable: fmt.Sprintf("showReviewDialogFor%s%s", app.Name, p.Username),
		},
		Status:           req.Status,
		StatusColor:      StatusColors[req.Status],
		StatusText:       StatusTexts[req.Status],
		StatusIncomplete: req.Status == StatusIncomplete,
		StatusReady:      req.Status == StatusReady,
		StatusSubmitted:  req.Status == StatusSubmitted,
		StatusInReview:   req.Status == StatusInReview,
		StatusApproved:   req.Status == StatusApproved,
		StatusReviewed:   req.Status == StatusReviewed,
		StatusRejected:   req.Status == StatusRejected,
		StatusArchived:   req.Status == StatusArchived,
		StatusCanceled:   req.Status == StatusCanceled,
		Link:             GetApplicationLink(p.ID, req, app),
		ID:               req.ID,
		Name:             name,
		Author:           p.Username,
		Reviewer:         reviewer,
		StatusIcon:       MakeStatusIcon(req.Status, 48),
	}
}

func GetApplicationLink(pid int64, req *queries.Request, app *queries.CharacterApplicationContent) string {
	// TODO: Update this to include reviewer's progress through the application
	if pid != req.PID {
		return routes.RequestFieldPath(req.ID, FieldName)
	}

	if !IsNameValid(app.Name) {
		return routes.RequestFieldPath(req.ID, FieldName)
	}
	if !IsGenderValid(app.Gender) {
		return routes.RequestFieldPath(req.ID, FieldGender)
	}
	if !IsShortDescriptionValid(app.ShortDescription) {
		return routes.RequestFieldPath(req.ID, FieldShortDescription)
	}
	if !IsDescriptionValid(app.Description) {
		return routes.RequestFieldPath(req.ID, FieldDescription)
	}
	if !IsBackstoryValid(app.Backstory) {
		return routes.RequestFieldPath(req.ID, FieldBackstory)
	}

	return routes.RequestFieldPath(req.ID, FieldBackstory)
}
