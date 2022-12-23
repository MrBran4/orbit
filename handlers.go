package orbit

import (
	"net/http"
)

// Similar to http.Handler but extended to handle route params.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, RouteParams, FromBodyable)
}
