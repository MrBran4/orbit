package orbit

import (
	"net/http"
)

// A Handler is similar to the stdlib'ss  http.Handler but extended to handle
// route params and an optional body.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, RouteParams, FromBodyable)
}

// HandlerFunc is Orbit's equivalent of net/http's HandlerFunc adapter.
// From the net/http docs:
//
// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request, RouteParams, FromBodyable)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, p RouteParams, b FromBodyable) {
	f(w, r, p, b)
}
