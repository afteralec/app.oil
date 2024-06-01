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
	VariableFulfill      = "showFulfillDialog"
)

const (
	TypePrimary     = "primary"
	TypeDestructive = "destructive"
)

type Definition struct {
	Header     string
	Text       template.HTML
	ButtonText string
	Path       string
	Variable   string
	Type       string
}

type DefinitionGroup struct {
	Submit       Definition
	Cancel       Definition
	PutInReview  Definition
	Approve      Definition
	FinishReview Definition
	Reject       Definition
	Fulfill      Definition
}

func (d *DefinitionGroup) Slice() []Definition {
	return []Definition{
		d.Submit,
		d.Cancel,
		d.PutInReview,
		d.Approve,
		d.FinishReview,
		d.Reject,
		d.Fulfill,
	}
}

func (d *DefinitionGroup) SetPath(rid int64) {
	path := route.RequestStatusPath(rid)
	d.Submit.Path = path
	d.Cancel.Path = path
	d.PutInReview.Path = path
	d.Approve.Path = path
	d.FinishReview.Path = path
	d.Fulfill.Path = path
}
