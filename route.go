package orbit

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"
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
	path     string       // The path to match on incoming requests. e.g. /a/b/{c}/d
	handler  Handler      // The handler which will actually deal with the request.
	params   RouteParams  // Route parameters that'll be passed to the handler (whose types must implement FromRequestable)
	bodyType FromBodyable // The type of the body (which will be nil if the handler doesn't care about the body or will decode its own)
	methods  []string     // The methods to match (e.g. get/put/patch). If it's empty, match all.
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
// ServeHTTP returns a few different error types:
//   - errRouteDoesNotMatch if the route simply doesn't match the request path.
//     In this case you should quietly continue trying against other handlers
//     in order until one does match.
//   - errMisconfigured if handling the request encouters something that looks
//     like it wasn't set up right (e.g. wrong number of args). This won't happen
//     intermittently - it'll either always work or never work. If you see this
//     it means you need to check how you're setting orbit up.
//   - Any other error - it'll bubble up errors returned by your FromRequest,
//     FromBody, or FromHeader funcs. Up to you how you want to handle these,
//     but best practice would probably be to return 5xx?
func (r *route) ServeHTTP(w http.ResponseWriter, req http.Request) error {

	// If this handler is set to match a specific method, check that.
	if len(r.methods) > 0 {
		if !contains(r.methods, strings.ToLower(req.Method)) {
			return errRouteDoesNotMatch("wrong http verb")
		}
	}

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

	// If the handler isn't expecting a decoded body, we can call it now.
	if r.bodyType == nil {
		r.handler.ServeHTTP(w, &req, *scopedParams, nil)
		return nil
	}

	// Try decoding the request body
	// Read the body and make 2 new readers from it, since reading once
	// consumes the body otherwise so it can't be re-read later.
	body, _ := io.ReadAll(req.Body)
	bReader1 := io.NopCloser(bytes.NewBuffer(body))
	bReader2 := io.NopCloser(bytes.NewBuffer(body))

	decodedBody, err := tryFromBody(r.bodyType, bReader1)
	if err != nil {
		return err
	}

	// Set the request's body back to the second reader so it's not empty anymore.
	req.Body = bReader2

	// Now call the handler, which will have all the params filled :)
	r.handler.ServeHTTP(w, &req, *scopedParams, decodedBody)

	return nil

}

// Kinda obvious but this checks if a slice of strings contains a string...
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
