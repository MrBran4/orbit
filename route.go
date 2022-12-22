package orbit

import (
	"net/http"
	"regexp"
)

// A route is an entry in a router.
// It defines the (possibly templated) path to match, and the handler to call.
//
// If the path contains {params} then the route will contain a pre-extracted
// list of those params along with their types (which is string by default, but
// can be anything), and a method for extracting the correct type (which is just
// to read the url by default, but can be any func.
//
// See the param docs for more info on how that works.
type route struct {
	// passed in during config:
	path    string      // The path to match on incoming requests. e.g. /a/b/{c}/d
	handler Handler     // The handler which will actually deal with the request.
	params  RouteParams // Route parameters that'll be passed to the handler (whose types must implement FromRequestable)
	// filters []FilterFunc // Request filters that can block execution if necessary (todo)

	// generated during config:
	orderedParamNames []string      // An ordered list of params in the path
	regex             regexp.Regexp // Precompiled regex for matching route params.
}

// Call bake when you're done configuring the routing tree. Call it only once.
// This 'precompiles' the handler by building the regex, param names etc.
func (r *route) bake() error {

	rx, opn, err := buildMatcherRegex(r.path)
	if err != nil {
		return err
	}

	r.regex = *rx
	r.orderedParamNames = opn

	return nil

}

// Calling ServeHTTP on a route causes it to handle the request if it matches.
// If the path doesn't match the path fed in, then the request won't be handled.
//
// ServeHTTP returns errRouteDoesNotMatch if the route doesn't match (in which
// case you should continue trying against other handlers), or any other error
// if this IS the right match but something goes wrong with handling the request
// (in which case you should not try any further handlers since this was the
// right one, it just failed somehow).
func (r *route) ServeHTTP(w http.ResponseWriter, req http.Request) error {

	// Try matching against the regex and extracting params.
	paramVals, err := tokenise(&r.regex, r.orderedParamNames, req.URL.Path)
	if err != nil {
		return err
	}

	// Build a param map populated with the ones from this request.
	// Note: If the params involve 'getting a user from the database based on
	//       an an ID provided in the request' etc. this is when that happens.
	scopedParams, err := r.params.newFromRequest(paramVals)
	if err != nil {
		return err
	}

	// Now call the handler, which will have all the params filled :)
	r.handler.ServeHTTP(w, &req, *scopedParams)

	return nil

}
