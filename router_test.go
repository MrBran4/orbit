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

func Benchmark_ServeHTTP_NoRouteParams_NoBody(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello",
		nil,
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/hello",
		handler,
		[]string{"POST"},
		nil,
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

func Benchmark_ServeHTTP_1RouteParam_NoBody(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/hello",
		nil,
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{foo}",
		handler,
		[]string{"POST"},
		RouteParams{
			"foo": testTypeString(""),
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

func Benchmark_ServeHTTP_10RouteParams_Nobody(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	// Test input
	req := httptest.NewRequest(
		http.MethodPost, // post
		"/a/b/one/two/three/four/five/six/seven/eight/nine/ten/c",
		nil,
	)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{p1}/{p2}/{p3}/{p4}/{p5}/{p6}/{p7}/{p8}/{p9}/{p10}/c",
		handler,
		[]string{"POST"},
		RouteParams{
			"p1":  testTypeString(""),
			"p2":  testTypeString(""),
			"p3":  testTypeString(""),
			"p4":  testTypeString(""),
			"p5":  testTypeString(""),
			"p6":  testTypeString(""),
			"p7":  testTypeString(""),
			"p8":  testTypeString(""),
			"p9":  testTypeString(""),
			"p10": testTypeString(""),
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

func Benchmark_ServeHTTP_NoRouteParams_Body(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/hello",
		handler,
		[]string{"POST"},
		nil,
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.NoError(b, err, "router bake failed")

	w := httptest.NewRecorder()

	// Reading an http.Request body consumes it.
	// Pre-allocate all the requests here to avoid artificially increasing the benchmark...
	var reqs []*http.Request
	for i := 0; i < b.N; i++ {
		reqs = append(reqs, httptest.NewRequest(
			http.MethodPost, // post
			"/a/b/hello",
			strings.NewReader("{}"),
		))
	}

	b.StartTimer()

	// Bench
	for i := 0; i < b.N; i++ {

		r.ServeHTTP(w, reqs[i])
	}

	// Check handler got called
	assert.Equal(b, b.N, calls)

}

func Benchmark_ServeHTTP_1RouteParam_Body(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{p1}",
		handler,
		[]string{"POST"},
		RouteParams{
			"p1": testTypeString(""),
		},
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.NoError(b, err, "router bake failed")

	w := httptest.NewRecorder()

	// Reading an http.Request body consumes it.
	// Pre-allocate all the requests here to avoid artificially increasing the benchmark...
	var reqs []*http.Request
	for i := 0; i < b.N; i++ {
		reqs = append(reqs, httptest.NewRequest(
			http.MethodPost, // post
			"/a/b/hello",
			strings.NewReader("{}"),
		))
	}

	b.StartTimer()

	// Bench
	for i := 0; i < b.N; i++ {

		r.ServeHTTP(w, reqs[i])
	}

	// Check handler got called
	assert.Equal(b, b.N, calls)

}

func Benchmark_ServeHTTP_10RouteParams_Body(b *testing.B) {

	// Stop bench timer while initialising
	b.StopTimer()

	calls := 0

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params RouteParams, body FromBodyable) {
		calls++
	})

	// Build a router, add the handler, bake
	r := NewRouter()
	r.Handle(
		"/a/b/{p1}/{p2}/{p3}/{p4}/{p5}/{p6}/{p7}/{p8}/{p9}/{p10}/c",
		handler,
		[]string{"POST"},
		RouteParams{
			"p1":  testTypeString(""),
			"p2":  testTypeString(""),
			"p3":  testTypeString(""),
			"p4":  testTypeString(""),
			"p5":  testTypeString(""),
			"p6":  testTypeString(""),
			"p7":  testTypeString(""),
			"p8":  testTypeString(""),
			"p9":  testTypeString(""),
			"p10": testTypeString(""),
		},
		testBodyableTypeStruct{},
	)

	err := r.Bake()
	assert.NoError(b, err, "router bake failed")

	w := httptest.NewRecorder()

	// Reading an http.Request body consumes it.
	// Pre-allocate all the requests here to avoid artificially increasing the benchmark...
	var reqs []*http.Request
	for i := 0; i < b.N; i++ {
		reqs = append(reqs, httptest.NewRequest(
			http.MethodPost, // post
			"/a/b/one/two/three/four/five/six/seven/eight/nine/ten/c",
			strings.NewReader("{}"),
		))
	}

	b.StartTimer()

	// Bench
	for i := 0; i < b.N; i++ {

		r.ServeHTTP(w, reqs[i])
	}

	// Check handler got called
	assert.Equal(b, b.N, calls)

}
