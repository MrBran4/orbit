package orbit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Route_ServeHTTP_Valid(t *testing.T) {

	// Flag - set true if the handler gets called (we want it to be called)
	handlerWasCalled := false

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello/d/123",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256
		}`))

	// Build & bake a route to test
	route := route{
		path: "/a/b/{foo}/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			// set flag true
			handlerWasCalled = true

			// Check everything's valid
			stringField, _ := params["foo"].(testTypeString)
			structField, _ := params["bar"].(testTypeStruct)
			bodyField, _ := body.(testBodyableTypeStruct)

			assert.Equal(t, "hello", string(stringField))
			assert.Equal(t, testTypeStruct{valuePassedIn: "123"}, structField)
			assert.Equal(t, testBodyableTypeStruct{FieldOne: "Hello World", FieldTwo: 256}, bodyField)
		}),
		params: RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		bodyType: testBodyableTypeStruct{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.NoError(t, err, "route bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	err = route.ServeHTTP(w, *req)
	assert.NoError(t, err, "ServeHTTP returned an err")

	// Check handler got called
	assert.True(t, handlerWasCalled, "looks like handler didn't get called")

}

func Test_Route_ServeHTTP_RouteDecodeFails(t *testing.T) {

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello/d/123",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256
		}`))

	// Build & bake a route to test
	route := route{
		path: "/a/b/{foo}/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			t.Fatalf("Handler was called when it shouldn't have been")
		}),
		params: RouteParams{
			"foo": testTypeBadIntWillFail(0),
			"bar": testTypeStruct{},
		},
		bodyType: testBodyableTypeStruct{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.NoError(t, err, "route bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	err = route.ServeHTTP(w, *req)
	assert.Error(t, err, "ServeHTTP should've halted and returned an err but didn't")

}

func Test_Route_ServeHTTP_BodyDecodeFails(t *testing.T) {

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello/d/123",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256
		}`))

	// Build & bake a route to test
	route := route{
		path: "/a/b/{foo}/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			t.Fatalf("Handler was called when it shouldn't have been")
		}),
		params: RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		bodyType: testBodyableTypeStructReturnsWrongType{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.NoError(t, err, "route bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	err = route.ServeHTTP(w, *req)
	assert.Error(t, err, "ServeHTTP should've halted and returned an err but didn't")

}

func Test_Route_ServeHTTP_WrongMethod(t *testing.T) {

	// Test input
	req := httptest.NewRequest(
		http.MethodGet, // get - handler's expecting post though
		"/a/b/hello/d/123",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256
		}`))

	// Build & bake a route to test
	route := route{
		path: "/a/b/{foo}/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			t.Fatalf("Handler was called when it shouldn't have been")
		}),
		params: RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		bodyType: testBodyableTypeStruct{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.NoError(t, err, "route bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	err = route.ServeHTTP(w, *req)
	assert.Error(t, err, "ServeHTTP should've halted and returned an err but didn't")

}

func Test_Route_ServeHTTP_WrongPath(t *testing.T) {

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello/d/123/e",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256
		}`))

	// Build & bake a route to test
	route := route{
		path: "/a/b/{foo}/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			t.Fatalf("Handler was called when it shouldn't have been")
		}),
		params: RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		bodyType: testBodyableTypeStruct{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.NoError(t, err, "route bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	err = route.ServeHTTP(w, *req)
	assert.Error(t, err, "ServeHTTP should've halted and returned an err but didn't")

}

func Test_Route_bake_WrongParamCount(t *testing.T) {
	route := route{
		path: "/a/b/{foo}/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			t.Fatalf("Handler was called when it shouldn't have been")
		}),
		params: RouteParams{
			"foo":               testTypeString(""),
			"bar":               testTypeStruct{},
			"mysteryThirdParam": testTypeStruct{},
		},
		bodyType: testBodyableTypeStruct{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.Error(t, err)
}

func Test_Route_bake_BadPath(t *testing.T) {
	route := route{
		path: "/a/b/{foo/d/{bar}",
		handler: HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
			t.Fatalf("Handler was called when it shouldn't have been")
		}),
		params: RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		bodyType: testBodyableTypeStruct{},
		methods:  []string{"POST"},
	}

	err := route.bake()
	assert.Error(t, err)
}
