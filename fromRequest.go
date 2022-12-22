package orbit

// For a type to be usable from a request it must implement FromRequestable.
//
// yourType.FromRequest(X) will be passed the string from the request param,
// and should set the element passed in to the value represented by that param.
// For example:
//
//	func (u User)FromRequest(uid) error {
//		// Set u to the user with ID uid,
//		// or return err if no user with that ID.
//	}
type FromRequestable interface {
	FromRequest(string) error
}
