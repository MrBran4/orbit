package orbit

import (
	"strconv"
)

// This file contains some helper types for decoding common things from params.

// BasicString is FromRequestable string type.
// You can use it to extract string values from url params.
type BasicString string

// FromRequest takes the raw URL param value (as a string) and returns a
// BasicString from it. It never returns an error.
func (x BasicString) FromRequest(param string) (any, error) {
	val := BasicString(param)
	return val, nil
}

// BasicInt is FromRequestable int type.
// You can use it to extract int values from url params.
type BasicInt int

// FromRequest takes the raw URL param value (as a string) and returns a
// BasicInt from it. If param is not a valid int, it returns an error.
func (x BasicInt) FromRequest(param string) (any, error) {
	intval, err := strconv.Atoi(param)
	if err != nil {
		return nil, err
	}

	result := BasicInt(intval)
	return result, nil
}
