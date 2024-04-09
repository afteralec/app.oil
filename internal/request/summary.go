package request

import (
	"context"
	"fmt"
	"html/template"

	"petrichormud.com/app/internal/player"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/dialog"
	"petrichormud.com/app/internal/request/field"
	"petrichormud.com/app/internal/route"
)

type SummaryForQueue struct {
	Dialogs         *dialog.DefinitionGroup
	StatusColor     string
	StatusText      string
	Title           string
	Link            string
	AuthorUsername  string
	ReviewerText    template.HTML
	StatusIcon      StatusIcon
	ID              int64
	PID             int64
	ShowPutInReview bool
}

type NewSummaryForQueueParams struct {
	Query               *query.Queries
	FieldMap            field.Map
	Request             *query.Request
	ReviewerPermissions *player.Permissions
	PlayerUsername      string
	ReviewerUsername    string
	PID                 int64
}

func NewSummaryForQueue(p NewSummaryForQueueParams) (SummaryForQueue, error) {
	// TODO: Move this into locals instead of updating the params
	player, err := p.Query.GetPlayer(context.Background(), p.Request.PID)
	if err != nil {
		return SummaryForQueue{}, err
	}
	p.PlayerUsername = player.Username
	if p.Request.RPID != 0 {
		reviewer, err := p.Query.GetPlayer(context.Background(), p.Request.RPID)
		if err != nil {
			return SummaryForQueue{}, err
		}
		p.ReviewerUsername = reviewer.Username
	}

	title := Title(p.Request.Type, p.FieldMap)

	reviewerText := ReviewerText(ReviewerTextParams{
		Request:          p.Request,
		ReviewerUsername: p.ReviewerUsername,
	})

	// TODO: Build a utility for this
	dialogs, ok := DialogsByType[p.Request.Type]
	if !ok {
		return SummaryForQueue{}, ErrNoDefinition
	}
	dialogs.SetPath(p.Request.ID)
	dialogs.PutInReview.Variable = fmt.Sprintf("showReviewDialogForRequest%d", p.Request.ID)
	showPutInReview := CanBePutInReview(
		CanBePutInReviewParams{
			Request:     p.Request,
			Permissions: p.ReviewerPermissions,
			PID:         p.PID,
		},
	)

	// TODO: Make this resilient to a request with an invalid status
	return SummaryForQueue{
		ID:              p.Request.ID,
		PID:             p.Request.PID,
		Title:           title,
		Link:            route.RequestPath(p.Request.ID),
		StatusIcon:      NewStatusIcon(StatusIconParams{Status: p.Request.Status, IconSize: 48, IncludeText: false}),
		StatusColor:     StatusColors[p.Request.Status],
		StatusText:      StatusTexts[p.Request.Status],
		ReviewerText:    reviewerText,
		Dialogs:         dialogs,
		AuthorUsername:  p.PlayerUsername,
		ShowPutInReview: showPutInReview,
	}, nil
}
