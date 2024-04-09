package request

import "errors"

const errNoDefinition string = "no definition with type"

var ErrNoDefinition error = errors.New(errNoDefinition)

var ErrInvalidType error = errors.New("invalid type")

var ErrInvalidInput error = errors.New("invalid input")
