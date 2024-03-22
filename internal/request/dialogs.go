package request

import (
	"html/template"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

const (
	BindCancelDialog       = "CancelDialog"
	BindSubmitDialog       = "SubmitDialog"
	BindPutInReviewDialog  = "PutInReviewDialog"
	BindApproveDialog      = "ApproveDialog"
	BindFinishReviewDialog = "FinishReviewDialog"
)

const (
	VariableCancelDialog      = "showCancelDialog"
	VariableSubmitDialog      = "showSubmitDialog"
	VariablePutInReviewDialog = "showPutInReviewDialog"
)

type Dialog struct {
	Header     string
	Text       template.HTML
	ButtonText string
	Path       string
	Variable   string
}

type Dialogs struct {
	Submit       Dialog
	Cancel       Dialog
	PutInReview  Dialog
	Approve      Dialog
	FinishReview Dialog
}

func (d *Dialogs) SetPath(rid int64) {
	path := route.RequestPath(rid)
	d.Submit.Path = path
	d.Cancel.Path = path
	d.PutInReview.Path = path
	d.Approve.Path = path
	d.FinishReview.Path = path
}

type BindDialogsParams struct {
	Request *query.Request
}

func BindDialogs(b fiber.Map, p BindDialogsParams) (fiber.Map, error) {
	def, ok := Definitions.Get(p.Request.Type)
	if !ok {
		return fiber.Map{}, ErrNoDefinition
	}

	dialogs := def.Dialogs()
	dialogs.SetPath(p.Request.ID)

	b[BindCancelDialog] = dialogs.Cancel
	b[BindSubmitDialog] = dialogs.Submit
	b[BindPutInReviewDialog] = dialogs.PutInReview
	b[BindApproveDialog] = dialogs.Approve
	b[BindFinishReviewDialog] = dialogs.FinishReview

	return b, nil
}
