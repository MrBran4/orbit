package orbit

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
	paths []route
}
