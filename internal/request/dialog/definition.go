package dialog

import (
	"html/template"

	"petrichormud.com/app/internal/route"
)

const (
	VariableCancel       = "showCancelDialog"
	VariableSubmit       = "showSubmitDialog"
	VariablePutInReview  = "showPutInReviewDialog"
	VariableApprove      = "showApproveDialog"
	VariableFinishReview = "showFinishReviewDialog"
	VariableReject       = "showRejectDialog"
)

type Definition struct {
	Header     string
	Text       template.HTML
	ButtonText string
	Path       string
	Variable   string
}

type DefinitionGroup struct {
	Submit       Definition
	Cancel       Definition
	PutInReview  Definition
	Approve      Definition
	FinishReview Definition
	Reject       Definition
}

func (d *DefinitionGroup) SetPath(rid int64) {
	path := route.RequestStatusPath(rid)
	d.Submit.Path = path
	d.Cancel.Path = path
	d.PutInReview.Path = path
	d.Approve.Path = path
	d.FinishReview.Path = path
}
