package orbit

import (
	"strconv"
)

// This file contains some helper types for decoding common things from params.

// A string from a request param as-is.
type BasicString string

func (x BasicString) FromRequest(param string) (any, error) {
	val := BasicString(param)
	return val, nil
}

// Extract an int from a request param.
type BasicInt int

func (x BasicInt) FromRequest(param string) (any, error) {
	intval, err := strconv.Atoi(param)
	if err != nil {
		return nil, err
	}

	result := BasicInt(intval)
	return result, nil
}
