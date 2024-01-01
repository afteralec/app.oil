package request

import (
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/permissions"
	"petrichormud.com/app/internal/queries"
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
	StatusSubmitted:  "fe:check-circle-o",
	StatusInReview:   "fe:question",
	StatusApproved:   "fe:check-circle",
	StatusReviewed:   "fe:warning",
	StatusRejected:   "fe:warning",
	StatusArchived:   "ic:round-lock",
	StatusCanceled:   "fe:outline-close",
}

var StatusColors map[string]string = map[string]string{
	StatusIncomplete: "text-gray-700",
	StatusReady:      "text-primary",
	StatusSubmitted:  "text-sky-700",
	StatusInReview:   "text-amber-700",
	StatusApproved:   "text-emerald-700",
	StatusReviewed:   "text-amber-700",
	StatusRejected:   "text-rose-700",
	StatusArchived:   "text-gray-700",
	StatusCanceled:   "text-gray-700",
}

// TODO: This can be shared across multiple packages
type StatusIcon struct {
	Size  string
	Color string
}

func IsStatusValid(status string) bool {
	_, ok := StatusTexts[status]
	return ok
}

func MakeStatusIcon(status string, size int64) StatusIcon {
	color, ok := StatusColors[status]
	if !ok {
		return StatusIcon{
			Size:  strconv.FormatInt(size, 10),
			Color: StatusColors[StatusIncomplete],
		}
	}

	return StatusIcon{
		Size:  strconv.FormatInt(size, 10),
		Color: color,
	}
}

func IsEditable(req *queries.Request) bool {
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

func IsStatusUpdateOK(req *queries.Request, perms permissions.PlayerGranted, pid int64, status string) bool {
	if status == StatusSubmitted {
		return req.PID == pid && req.Status == StatusReady
	}

	if status == StatusCanceled {
		return req.PID == pid
	}

	if status == StatusInReview {
		if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
			return false
		}
		return req.PID != pid && req.Status == StatusSubmitted
	}

	if status == StatusReviewed {
		if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
			return false
		}
		return req.PID != pid && req.Status == StatusInReview
	}

	if status == StatusApproved {
		if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
			return false
		}
		return req.PID != pid && req.Status == StatusInReview
	}

	if status == StatusRejected {
		if !perms.Permissions[permissions.PlayerReviewCharacterApplicationsName] {
			return false
		}
		return req.PID != pid && req.Status == StatusInReview
	}

	return false
}

func BindStatuses(b fiber.Map, req *queries.Request) fiber.Map {
	b["StatusIncomplete"] = StatusIncomplete
	b["StatusReady"] = StatusReady
	b["StatusSubmitted"] = StatusSubmitted
	b["StatusInReview"] = StatusInReview
	b["StatusApproved"] = StatusApproved
	b["StatusReviewed"] = StatusReviewed
	b["StatusRejected"] = StatusRejected
	b["StatusArchived"] = StatusArchived
	b["StatusCanceled"] = StatusCanceled

	b["StatusIsIncomplete"] = req.Status == StatusIncomplete
	b["StatusIsReady"] = req.Status == StatusReady
	b["StatusIsSubmitted"] = req.Status == StatusSubmitted
	b["StatusIsInReview"] = req.Status == StatusInReview
	b["StatusIsApproved"] = req.Status == StatusApproved
	b["StatusIsReviewed"] = req.Status == StatusReviewed
	b["StatusIsRejected"] = req.Status == StatusRejected
	b["StatusIsArchived"] = req.Status == StatusArchived
	b["StatusIsCanceled"] = req.Status == StatusCanceled

	b["StatusText"] = StatusTexts[req.Status]

	b["StatusColor"] = StatusColors[req.Status]

	return b
}
