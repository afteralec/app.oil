package route

import (
	"fmt"
	"strings"
)

const (
	Requests                    = "/requests"
	RequestPathParam            = "/requests/:id"
	RequestFieldPathParam       = "/requests/:id/:field"
	RequestFieldStatusPathParam = "/requests/:id/:field/status"
	RequestStatusPathParam      = "/requests/:id/status"
)

const CreateRequestCommentPathParam string = "/request/:id/comment/:field"

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

func RequestStatusPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/status", Requests, id)
	return b.String()
}

func CreateRequestCommentPath(rid int64, field string) string {
	return fmt.Sprintf("/request/%d/comment/%s", rid, field)
}
