package orbit

import "io"

// For a type to be usable from a request it must implement FromRequestable.
//
// yourType.FromRequest(X) will be called by orbit when a request is being handled
// where one of the parameters is expected to be of that type.

// For example if a request to /user/{uid} is called with uid=5, and you
// handler has specified the {uid} param to be of type yourType, then then orbit
// will call yourType.FromRequest("5"), and expects it to return (yourType, nil).
//
// Your type's FromRequest func should return the value represented by that param.
// So in the example above you should return a yourType struct representing the
// user with uid 5 (or return an error if you can't do that, e.g. if there
// isn't a user with uid 5.
//
// For example:
//
//	func (u User)FromRequest(uid) (any, error) {
//		// Parse uid to an int,
//		// Get the user with that uid from the database
//		// Build a User struct from that data and return it.
//	}
type FromRequestable interface {
	FromRequest(string) (any, error)
}

// For a type to be decodable from a request body, it must implement FromBodyable.
//
// yourType.FromBody(X) will be called by orbit when a request is being handled
// whose body is expected to be of that type.
//
// X will be an io.ReadCloser for the request body (i.e. reading from the reader
// is like reading from a http.Request.Body), and you should close it.
//
// For example, for some fictional 'submit a new todo' endpoint:
//
//	func (u Todo)FromBody(body) (any, error) {
//		// Poor error handling for sake of example
//		var todo Todo
//		_ := json.NewDecoder(body).Decode(&todo);
//		return todo, nil
//	}
type FromBodyable interface {
	FromBody(io.ReadCloser) (any, error)
}
