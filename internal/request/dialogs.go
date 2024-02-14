package request

import (
	"html/template"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/route"
)

const (
	BindCancelDialog      = "CancelDialog"
	BindSubmitDialog      = "SubmitDialog"
	BindPutInReviewDialog = "PutInReviewDialog"
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
	Submit      Dialog
	Cancel      Dialog
	PutInReview Dialog
}

var BindDialogsByType map[string]Dialogs = map[string]Dialogs{
	TypeCharacterApplication: {
		Submit: Dialog{
			Header:     "Submit This Application?",
			Text:       "Once your character application is put in review, this cannot be undone.",
			ButtonText: "Submit This Application",
		},
		Cancel: Dialog{
			Header:     "Cancel This Application?",
			Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
			ButtonText: "Cancel This Application",
		},
		PutInReview: Dialog{
			Header:     "Put This Application In Review?",
			Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
			ButtonText: "I'm Ready to Review This Application",
		},
	},
}

type BindDialogsParams struct {
	Request *query.Request
}

func BindDialogs(b fiber.Map, p BindDialogsParams) fiber.Map {
	bindDialogs, ok := BindDialogsByType[p.Request.Type]
	if !ok {
		return b
	}

	b[BindCancelDialog] = Dialog{
		Header:     bindDialogs.Cancel.Header,
		Text:       bindDialogs.Cancel.Text,
		ButtonText: bindDialogs.Cancel.ButtonText,
		Path:       route.RequestPath(p.Request.ID),
		Variable:   VariableCancelDialog,
	}

	b[BindSubmitDialog] = Dialog{
		Header:     bindDialogs.Submit.Header,
		Text:       bindDialogs.Submit.Text,
		ButtonText: bindDialogs.Submit.ButtonText,
		Path:       route.RequestStatusPath(p.Request.ID),
		Variable:   VariableSubmitDialog,
	}

	b[BindPutInReviewDialog] = Dialog{
		Header:     bindDialogs.PutInReview.Header,
		Text:       bindDialogs.PutInReview.Text,
		ButtonText: bindDialogs.PutInReview.ButtonText,
		Path:       route.RequestStatusPath(p.Request.ID),
		Variable:   VariablePutInReviewDialog,
	}

	return b
}
