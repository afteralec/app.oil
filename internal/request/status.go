package request

import (
	"context"
	"errors"
	"html/template"

	fiber "github.com/gofiber/fiber/v2"

	"petrichormud.com/app/internal/query"
)

const (
	StatusIncomplete = "Incomplete"
	StatusReady      = "Ready"
	StatusSubmitted  = "Submitted"
	StatusInReview   = "InReview"
	StatusApproved   = "Approved"
	StatusReviewed   = "Reviewed"
	StatusRejected   = "Rejected"
	StatusArchived   = "Archived"
	StatusCanceled   = "Canceled"
)

var StatusTexts map[string]string = map[string]string{
	StatusIncomplete: "Incomplete",
	StatusReady:      "Ready",
	StatusSubmitted:  "Submitted",
	StatusInReview:   "In Review",
	StatusApproved:   "Approved",
	StatusReviewed:   "Reviewed",
	StatusRejected:   "Rejected",
	StatusArchived:   "Archived",
	StatusCanceled:   "Canceled",
}

var StatusIcons map[string]string = map[string]string{
	StatusIncomplete: "ph:dots-three-outline-fill",
	StatusReady:      "fe:check",
	StatusSubmitted:  "gg:check-o",
	StatusInReview:   "fe:question",
	StatusApproved:   "gg:check-o",
	StatusReviewed:   "fe:warning",
	StatusRejected:   "fe:warning",
	StatusArchived:   "ic:round-lock",
	StatusCanceled:   "fe:outline-close",
}

var StatusColors map[string]string = map[string]string{
	StatusIncomplete: "text-incomplete",
	StatusReady:      "text-ready",
	StatusSubmitted:  "text-submitted",
	StatusInReview:   "text-review",
	StatusApproved:   "text-approved",
	StatusReviewed:   "text-reviewed",
	StatusRejected:   "text-rejected",
	StatusArchived:   "text-archived",
	StatusCanceled:   "text-canceled",
}

type StatusIcon struct {
	Icon  template.URL
	Size  string
	Color string
	Text  string
}

func IsStatusValid(status string) bool {
	_, ok := StatusTexts[status]
	return ok
}

type MakeStatusIconParams struct {
	Status      string
	Size        string
	IncludeText bool
}

func MakeStatusIcon(p MakeStatusIconParams) StatusIcon {
	icon, ok := StatusIcons[p.Status]
	if !ok {
		return MakeDefaultStatusIcon(p.Size, p.IncludeText)
	}

	color, ok := StatusColors[p.Status]
	if !ok {
		return MakeDefaultStatusIcon(p.Size, p.IncludeText)
	}

	result := StatusIcon{
		Icon:  template.URL(icon),
		Color: color,
		Size:  p.Size,
	}

	if p.IncludeText {
		text, ok := StatusTexts[p.Status]
		if !ok {
			return MakeDefaultStatusIcon(p.Size, p.IncludeText)
		}

		result.Text = text
	}

	return result
}

func MakeDefaultStatusIcon(size string, includeText bool) StatusIcon {
	result := StatusIcon{
		Icon:  template.URL(StatusIcons[StatusIncomplete]),
		Color: StatusColors[StatusIncomplete],
		Size:  size,
	}

	if includeText {
		result.Text = StatusTexts[StatusIncomplete]
	}

	return result
}

func IsEditable(req *query.Request) bool {
	if req.Status == StatusSubmitted {
		return false
	}

	if req.Status == StatusInReview {
		return false
	}

	if req.Status == StatusReviewed {
		return false
	}

	if req.Status == StatusApproved {
		return false
	}

	if req.Status == StatusRejected {
		return false
	}

	if req.Status == StatusCanceled {
		return false
	}

	if req.Status == StatusArchived {
		return false
	}

	return true
}

func BindStatus(b fiber.Map, req *query.Request) fiber.Map {
	b["StatusIsIncomplete"] = req.Status == StatusIncomplete
	b["StatusIsReady"] = req.Status == StatusReady
	b["StatusIsSubmitted"] = req.Status == StatusSubmitted
	b["StatusIsInReview"] = req.Status == StatusInReview
	b["StatusIsApproved"] = req.Status == StatusApproved
	b["StatusIsReviewed"] = req.Status == StatusReviewed
	b["StatusIsRejected"] = req.Status == StatusRejected
	b["StatusIsArchived"] = req.Status == StatusArchived
	b["StatusIsCanceled"] = req.Status == StatusCanceled

	return b
}

type UpdateStatusParams struct {
	Status string
	PID    int64
	RID    int64
}

var ErrInvalidStatus error = errors.New("invalid status")

func UpdateStatus(q *query.Queries, p UpdateStatusParams) error {
	if !IsStatusValid(p.Status) {
		return ErrInvalidStatus
	}

	if err := q.UpdateRequestStatus(context.Background(), query.UpdateRequestStatusParams{
		ID:     p.RID,
		Status: p.Status,
	}); err != nil {
		return err
	}

	if err := q.CreateHistoryForRequestStatusChange(context.Background(), query.CreateHistoryForRequestStatusChangeParams{
		RID: p.RID,
		PID: p.PID,
	}); err != nil {
		return err
	}

	return nil
}
