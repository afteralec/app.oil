package request

import (
	"html/template"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/routes"
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

type BindDialog struct {
	Header     string
	Text       template.HTML
	ButtonText string
	Path       string
	Variable   string
}

type BindDialogsByKind struct {
	Submit      BindDialog
	Cancel      BindDialog
	PutInReview BindDialog
}

var BindDialogsByType map[string]BindDialogsByKind = map[string]BindDialogsByKind{
	TypeCharacterApplication: {
		Submit: BindDialog{
			Header:     "Submit This Application?",
			Text:       "Once your character application is put in review, this cannot be undone.",
			ButtonText: "Submit This Application",
		},
		Cancel: BindDialog{
			Header:     "Cancel This Application?",
			Text:       "Once you've canceled this application, it cannot be undone. If you want to apply with this character again in the future, you'll need to create a new application.",
			ButtonText: "Cancel This Application",
		},
		PutInReview: BindDialog{
			Header:     "Put This Application In Review?",
			Text:       template.HTML("Once you put this application in review, <span class=\"font-semibold\">you must review it within 24 hours</span>. After picking up this application, you'll be the only reviewer able to review it."),
			ButtonText: "I'm Ready to Review This Application",
		},
	},
}

type BindDialogsParams struct {
	Request *queries.Request
}

func BindDialogs(b fiber.Map, p BindDialogsParams) fiber.Map {
	bindDialogs, ok := BindDialogsByType[p.Request.Type]
	if !ok {
		return b
	}

	b[BindCancelDialog] = BindDialog{
		Header:     bindDialogs.Cancel.Header,
		Text:       bindDialogs.Cancel.Text,
		ButtonText: bindDialogs.Cancel.ButtonText,
		Path:       routes.RequestPath(p.Request.ID),
		Variable:   VariableCancelDialog,
	}

	b[BindSubmitDialog] = BindDialog{
		Header:     bindDialogs.Submit.Header,
		Text:       bindDialogs.Submit.Text,
		ButtonText: bindDialogs.Submit.ButtonText,
		Path:       routes.RequestPath(p.Request.ID),
		Variable:   VariableSubmitDialog,
	}

	b[BindPutInReviewDialog] = BindDialog{
		Header:     bindDialogs.PutInReview.Header,
		Text:       bindDialogs.PutInReview.Text,
		ButtonText: bindDialogs.PutInReview.ButtonText,
		Path:       routes.RequestPath(p.Request.ID),
		Variable:   VariablePutInReviewDialog,
	}
	return b
}
