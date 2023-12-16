package request

import (
	fiber "github.com/gofiber/fiber/v2"

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

func BindStatuses(b fiber.Map, req *queries.Request) fiber.Map {
	b["StatusIncomplete"] = req.Status == StatusIncomplete
	b["StatusReady"] = req.Status == StatusReady
	b["StatusSubmitted"] = req.Status == StatusSubmitted
	b["StatusInReview"] = req.Status == StatusInReview
	b["StatusApproved"] = req.Status == StatusApproved
	b["StatusReviewed"] = req.Status == StatusReviewed
	b["StatusRejected"] = req.Status == StatusRejected
	b["StatusArchived"] = req.Status == StatusArchived
	b["StatusCanceled"] = req.Status == StatusCanceled
	return b
}
