package orbit

import (
	"fmt"
	"reflect"
)

// RouteParams are the values extracted from URL parameters (e.g. /from/{this}/)
//
// It's a map of[string]FromRequestable where the key is the parameter name from
// the url, and the value is the object decoded from that.
//
// For a route like /users/{user}/photos/{photo}, your RouteParms would be like:
//
//	RouteParams{
//		"user": YourUserType{},
//		"photo": YourPhotoType{},
//	}
//
// When Orbit calls your handler, it'll pass the same RouteParams to you, filled
// with the data it decoded from the url.
type RouteParams map[string]FromRequestable

// newFromRequest takes a request path (e.g. /a/b/{c}/d) and returns a copy
// with all arguments populated from that request.
//
// To be successful, *all* fields must populate correctly. If any fields fail
// to populate, then an error is returned.
func (params RouteParams) newFromRequest(tokens map[string]string) (*RouteParams, error) {

	filled := make(RouteParams)
	rvfilled := reflect.ValueOf(filled)

	for key, el := range params {
		result, err := el.FromRequest(tokens[key])
		if err != nil {
			return nil, errCoudlntGetParams{paramName: key, err: err}
		}

		// use reflection to check the type returned is correct
		if reflect.TypeOf(result) != reflect.TypeOf(el) {
			return nil, errMisconfigured(fmt.Sprintf("%s's FromRequest method returned unexpected type (want %s got %s)", key, reflect.TypeOf(el), reflect.TypeOf(result)))
		}
		rvfilled.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(result))
	}

	return &filled, nil

}
