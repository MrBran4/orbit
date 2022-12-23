package orbit

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
