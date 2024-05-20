package route

import (
	"fmt"
	"strings"
)

// TODO: Generate these with a param from the string blocks below
const (
	Requests                           = "/requests"
	RequestPathParam                   = "/requests/:id"
	RequestFieldPathParam              = "/requests/:rid/fields/:rfid"
	RequestSubfieldsPathParam          = "/requests/:rid/fields/:rfid/subfields"
	RequestSubfieldPathParam           = "/requests/:rid/fields/:rfid/subfields/:id"
	RequestFieldTypePathParam          = "/requests/:id/:field"
	RequestFieldStatusPathParam        = "/requests/:id/:field/status"
	RequestChangeRequestPathParam      = "/requests/changes/:id"
	RequestChangeRequestFieldPathParam = "/requests/:id/:field/changes"
	RequestStatusPathParam             = "/requests/:id/status"
)

const (
	RequestFields         = "fields"
	RequestSubfields      = "subfields"
	RequestChangeRequests = "changes"
)

func RequestPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d", Requests, id)
	return b.String()
}

func RequestFieldTypePath(id int64, field string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s", Requests, id, field)
	return b.String()
}

func RequestFieldStatusPath(id int64, field string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s/status", Requests, id, field)
	return b.String()
}

func RequestFieldSubfieldsPath(rid, id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s/%d/%s", Requests, rid, RequestFields, id, RequestSubfields)
	return b.String()
}

func RequestFieldSubfieldPath(rid, rfid, id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s/%d/%s/%d", Requests, rid, RequestFields, rfid, RequestSubfields, id)
	return b.String()
}

func RequestChangeRequestPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%s/%d", Requests, RequestChangeRequests, id)
	return b.String()
}

func RequestChangeRequestFieldPath(id int64, field string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/%s/%s", Requests, id, field, RequestChangeRequests)
	return b.String()
}

func RequestStatusPath(id int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/%d/status", Requests, id)
	return b.String()
}
