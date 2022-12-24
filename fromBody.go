package orbit

import (
	"fmt"
	"io"
	"reflect"
)

// The FromBodyable interface allows Orbit to resolve your type from a request
// body (by calling your type\s FromBody function, and passing it the body).
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

func tryFromBody(bodyType FromBodyable, body io.ReadCloser) (FromBodyable, error) {

	var resultMap map[string]FromBodyable = map[string]FromBodyable{"result": bodyType}
	rResultMap := reflect.ValueOf(resultMap)

	// Try decoding the body (as 'any' type) by calling the type's FromBody.
	decodedBodyAsAny, err := bodyType.FromBody(body)
	if err != nil {
		return nil, err
	}

	// Decoded body currently has type 'any'. Use reflection to check it's
	// actually of the correct type...
	rExpectedType := reflect.ValueOf(bodyType)
	rDecodedBodyAsAny := reflect.ValueOf(decodedBodyAsAny)
	if rDecodedBodyAsAny.Type() != rExpectedType.Type() {
		return nil, errMisconfigured(fmt.Sprintf("FromBody method returned unexpected type (want %s got %s)", rExpectedType.Type(), rDecodedBodyAsAny.Type()))
	}
	rResultMap.SetMapIndex(reflect.ValueOf("result"), rDecodedBodyAsAny)

	return resultMap["result"], nil

}
