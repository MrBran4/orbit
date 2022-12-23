package orbit

import (
	"fmt"
	"log"
	"net/http"
)

// A Router routes an incoming request to handler, based on its path and its method.
//
// Paths can contain parameters in squirly braces. For example: /a/b/{c}/d matches
// requests to /a/b/ANY_STRING_HERE/d. The value of parameter {c} would then be
// available in the handler.
//
// It's a `net/http` compliant handler, so you can call .ServeHTTP one it for a
// one-off request or, more usefully, you can use it in .ListenAndServe.
//
// Build the router with .Handler or .Subrouter calls
type Router struct {
	routes []route
}

// Create a new Orbit router, off which you can hang your handlers.
func NewRouter() Router {
	return Router{}
}

// Add a new handler to the router.
//   - path is your /parameterised/{path}/to/{match}
//   - handler is your handler to call if the request is valid
//   - methods is a []string of HTTP verbs to match, or nil to match all of them
//   - routeParamTypes a map of paramName/TypeOfParam that the handler expects, or nil if the route has no params.
//   - bodyType is the type the body should be decoded to, or nil for orbit to skip decoding that.
//
// The handler will be called if:
//   - the path matches path
//   - the methods match
//   - all params in the request successfully resolve to the types specified in routeParamTypes
//   - the body successfully resolves to bodyType (if bodyType isn't nil)
func (r *Router) Handle(
	path string,
	handler Handler,
	methods []string,
	routeParamTypes RouteParams,
	bodyType FromBodyable,
) {

	// if routeParamTypes is nil, set it to an empty map.
	paramTypes := routeParamTypes
	if routeParamTypes == nil {
		paramTypes = make(RouteParams)
	}

	// append it to the route
	r.routes = append(r.routes, route{
		path:     path,
		handler:  handler,
		methods:  methods,
		params:   paramTypes,
		bodyType: bodyType,
	})

}

// Bake prepares the router to be used.
// It precompiles your routes' regexes and checks the params match up.
//
// Call Bake exactly once, after you have added all of your routes and before you
// start using the router.
func (r *Router) Bake() error {

	for i := 0; i < len(r.routes); i++ {
		if err := r.routes[i].bake(); err != nil {
			return errMisconfigured(fmt.Sprintf("couldn't bake handler '%s': %s", r.routes[i].path, err.Error()))
		}
	}

	return nil

}

// Handle an incoming HTTP request
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Try every handler until one handles it.
	for _, route := range router.routes {
		err := route.ServeHTTP(w, *r)

		// If this handler successfully handled the route, we can stop searching
		if err == nil {
			return
		}

		// If the error is errRouteDoesNotMatch, try the next handler
		if _, ok := err.(*errRouteDoesNotMatch); ok {
			continue
		}

		// If the err is any other type, stop processing handlers and log it
		log.Printf("orbit encountered an error handling '%s': %s\n", r.URL.Path, err.Error())
		return
	}

}
