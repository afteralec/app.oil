package request

import (
	"errors"
)

// TODO: This can likely be moved or removed
var ErrInvalidInput error = errors.New("field value didn't pass validation")
