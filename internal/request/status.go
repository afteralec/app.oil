package request

import "petrichormud.com/app/internal/queries"

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
