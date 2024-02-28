package request

import (
	"errors"
)

const (
	FieldName             string = "name"
	FieldGender           string = "gender"
	FieldShortDescription string = "sdesc"
	FieldDescription      string = "desc"
	FieldBackstory        string = "backstory"
)

// TODO: This can likely be moved or removed
var ErrInvalidInput error = errors.New("field value didn't pass validation")
