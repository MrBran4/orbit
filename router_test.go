package orbit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Router_E2E_Valid(t *testing.T) {

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

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		// set flag true
		handlerWasCalled = true

		// Check everything's valid
		stringField, _ := params["foo"].(testTypeString)
		structField, _ := params["bar"].(testTypeStruct)
		bodyField, _ := body.(testBodyableTypeStruct)

		assert.Equal(t, "hello", string(stringField))
		assert.Equal(t, testTypeStruct{valuePassedIn: "123"}, structField)
		assert.Equal(t, testBodyableTypeStruct{FieldOne: "Hello World", FieldTwo: 256}, bodyField)
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{foo}/d/{bar}",
		handler,
		[]string{"POST"},
		RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.NoError(t, err, "router bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check handler got called
	assert.True(t, handlerWasCalled, "looks like handler didn't get called")

}

func Test_Router_E2E_NoBody(t *testing.T) {

	// Flag - set true if the handler gets called (we want it to be called)
	handlerWasCalled := false

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello/d/123",
		nil,
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		// set flag true
		handlerWasCalled = true

		// Check everything's valid
		stringField, _ := params["foo"].(testTypeString)
		structField, _ := params["bar"].(testTypeStruct)

		assert.Equal(t, "hello", string(stringField))
		assert.Equal(t, testTypeStruct{valuePassedIn: "123"}, structField)
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{foo}/d/{bar}",
		handler,
		[]string{"POST"},
		RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		nil,
	)

	err := r.Bake()
	assert.NoError(t, err, "router bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check handler got called
	assert.True(t, handlerWasCalled, "looks like handler didn't get called")

}

func Test_Router_E2E_NoParams(t *testing.T) {

	// Flag - set true if the handler gets called (we want it to be called)
	handlerWasCalled := false

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/c/d/e",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256
		}`))

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		// set flag true
		handlerWasCalled = true

		// Check everything's valid
		bodyField, _ := body.(testBodyableTypeStruct)
		assert.Equal(t, testBodyableTypeStruct{FieldOne: "Hello World", FieldTwo: 256}, bodyField)
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/c/d/e",
		handler,
		[]string{"POST"},
		nil,
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.NoError(t, err, "router bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check handler got called
	assert.True(t, handlerWasCalled, "looks like handler didn't get called")

}

func Test_Router_E2E_NoHandler(t *testing.T) {

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/c/d/e/f",
		nil,
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		t.Fatalf("handler was called when it shouldn't have been")
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/c/d/e",
		handler,
		[]string{"POST"},
		nil,
		nil,
	)

	err := r.Bake()
	assert.NoError(t, err, "router bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check Orbit returned 404
	assert.Equal(t, 404, w.Code)

}

func Test_Router_E2E_DecodeFails(t *testing.T) {

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/c/d/e",
		strings.NewReader(`{
			"field_one": "Hello World",
			"field_two": 256 i am invalid json!
		}`),
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		t.Fatalf("handler was called when it shouldn't have been")
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/c/d/e",
		handler,
		[]string{"POST"},
		nil,
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.NoError(t, err, "router bake failed")

	// Handle the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check Orbit returned 503
	assert.Equal(t, 503, w.Code)

}

func Test_Router_E2E_Misconfiguration(t *testing.T) {

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		t.Fatalf("handler was called when it shouldn't have been")
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{c/d/{e}", // <-- Unbalanced braces
		handler,
		[]string{"POST"},
		nil,
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.Error(t, err, "router bake failed")

}

func Benchmark_Router_ServeHTTP_2Params(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello/d/123",
		nil,
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{foo}/d/{bar}",
		handler,
		[]string{"POST"},
		RouteParams{
			"foo": testTypeString(""),
			"bar": testTypeStruct{},
		},
		nil,
	)

	err := r.Bake()
	assert.NoError(b, err, "router bake failed")

	w := httptest.NewRecorder()

	b.StartTimer()

	// Bench
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}

	// Check handler got called
	assert.Equal(b, b.N, calls)

}
