package filter

import "errors"

func NewFilterError(msg string) (err error) {
	return errors.New(msg)
}

var ErrorFilterInvalidOperation = NewFilterError("invalid operation")
var ErrorFilterUnableToParseQuery = NewFilterError("unable to parse query")
