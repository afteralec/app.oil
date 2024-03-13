package route

import (
	"fmt"
	"strings"
)

const (
	Requests                           = "/requests"
	RequestPathParam                   = "/requests/:id"
	RequestFieldPathParam              = "/requests/:id/:field"
	RequestFieldStatusPathParam        = "/requests/:id/:field/status"
	RequestChangeRequestPathParam      = "/requests/changes/:id"
	RequestChangeRequestFieldPathParam = "/requests/:id/:field/changes"
	RequestStatusPathParam             = "/requests/:id/status"
)

const ChangeRequests = "changes"

func RequestPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d", Requests, id)
	return b.String()
}

func RequestFieldPath(id int64, field string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s", Requests, id, field)
	return b.String()
}

func RequestFieldStatusPath(id int64, field string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s/status", Requests, id, field)
	return b.String()
}

func RequestChangeRequestPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%s/%d", Requests, ChangeRequests, id)
	return b.String()
}

func RequestChangeRequestFieldPath(id int64, field string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s/%s", Requests, id, field, ChangeRequests)
	return b.String()
}

func RequestStatusPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/status", Requests, id)
	return b.String()
}
