package orbit

import "fmt"

type errRouteDoesNotMatch string

func (e errRouteDoesNotMatch) Error() string {
	return fmt.Sprintf("route doesn't match (%s)", string(e))
}

type errMisconfigured string

func (e errMisconfigured) Error() string {
	return fmt.Sprintf("orbit may be misconfigured: %s", string(e))
}

type errCoudlntGetParams struct {
	paramName string
	err       error
}

func (e errCoudlntGetParams) Error() string {
	return fmt.Sprintf("couldn't get %s from request (%s)", e.paramName, e.err.Error())
}
