package request

import (
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/internal/request/field"
)

func IsEditable(pid int64, req *query.Request, fd field.Field) bool {
	if fd.ForReviewer() {
		return req.Status == StatusInReview && pid == req.RPID
	}

	if fd.ForPlayer() {
		return pid == req.PID && (req.Status == StatusIncomplete ||
			req.Status == StatusReady ||
			req.Status == StatusReviewed)
	}

	return false
}
