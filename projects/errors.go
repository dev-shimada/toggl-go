package projects

import (
	"errors"
)

var (
	ErrorStatusNotOK       = errors.New("error response status code")
	ErrorRequiredParameter = errors.New("required parameter is missing")
)
