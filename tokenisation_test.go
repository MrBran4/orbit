package orbit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getPositionsOfSquirlies(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []int
		wantErr bool
	}{
		{name: "valid_1", path: "foo/bar/{param1}/test", want: []int{8, 16}, wantErr: false},
		{name: "valid_2", path: "foo/bar/{param1}/test/{param2}/baz", want: []int{8, 16, 22, 30}, wantErr: false},
		{name: "valid_none", path: "foo/bar/test/", want: []int{}, wantErr: false},
		{name: "valid_start", path: "{param1}/foo/bar/test/", want: []int{0, 8}, wantErr: false},
		{name: "valid_end", path: "foo/bar/{param1}", want: []int{8, 16}, wantErr: false},
		{name: "unbalanced", path: "foo/bar/{param1/test/{param2}/baz", want: nil, wantErr: true},
		{name: "unclosed", path: "foo/bar/{param1}/test/{param2/baz", want: nil, wantErr: true},
		{name: "unopened", path: "foo/bar/param1}/test/{param2}/baz", want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPositionsOfSquirlies(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPositionsOfSquirlies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// Checks the regex building function works.
//
// DOES NOT check whether that regex is actually what we want, just that it's a valid regex.
// There are integration tests that cover whether the regex does the right job.
func Test_buildMatcherRegex(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantNames []string
		wantErr   bool
	}{
		{name: "valid_1", path: "/aaa/bbb/{foo}/ccc/ddd", wantNames: []string{"foo"}, wantErr: false},
		{name: "valid_2", path: "/aaa/bbb/{foo}/ccc/{bar}/ddd", wantNames: []string{"foo", "bar"}, wantErr: false},
		{name: "valid_none", path: "/aaa/bbb/ccc/ddd", wantNames: []string{}, wantErr: false},
		{name: "valid_start", path: "{foo}/aaa/bbb/ccc/{bar}/ddd", wantNames: []string{"foo", "bar"}, wantErr: false},
		{name: "valid_end", path: "{foo}/aaa/bbb/ccc/ddd/{bar}", wantNames: []string{"foo", "bar"}, wantErr: false},
		{name: "valid_touching", path: "/aaa/bbb/{foo}/{bar}/ccc/ddd/", wantNames: []string{"foo", "bar"}, wantErr: false},
		{name: "valid_really_touching", path: "/aaa/bbb/{foo}{bar}/ccc/ddd/", wantNames: []string{"foo", "bar"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotNames, err := buildMatcherRegex(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Test_buildMatcherRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantNames, gotNames)
		})
	}
}

func Test_IntegrationTest_tokenisation(t *testing.T) {
	tests := []struct {
		name         string
		path         string            // Template
		reqPath      string            // 'Real request' path
		wantParamMap map[string]string // expected extraction
		wantErr      bool              // want an error using regex?
	}{
		{name: "valid_1",
			path:    "/aaa/bbb/{foo}/ccc/ddd",
			reqPath: "/aaa/bbb/hello/ccc/ddd",
			wantParamMap: map[string]string{
				"foo": "hello",
			},
			wantErr: false,
		},
		{name: "valid_2",
			path:    "/aaa/bbb/{foo}/ccc/{bar}/ddd",
			reqPath: "/aaa/bbb/hello/ccc/world/ddd",
			wantParamMap: map[string]string{
				"foo": "hello",
				"bar": "world",
			},
			wantErr: false,
		},
		{name: "valid_none",
			path:         "/aaa/bbb/ccc/ddd",
			reqPath:      "/aaa/bbb/ccc/ddd",
			wantParamMap: map[string]string{},
			wantErr:      false,
		},
		{name: "valid_start",
			path:    "{foo}/aaa/bbb/ccc/{bar}/ddd",
			reqPath: "hello/aaa/bbb/ccc/world/ddd",
			wantParamMap: map[string]string{
				"foo": "hello",
				"bar": "world",
			},
			wantErr: false,
		},
		{name: "valid_end",
			path:    "{foo}/aaa/bbb/ccc/ddd/{bar}",
			reqPath: "hello/aaa/bbb/ccc/ddd/world",
			wantParamMap: map[string]string{
				"foo": "hello",
				"bar": "world",
			},
			wantErr: false,
		},
		{name: "valid_touching",
			path:    "/aaa/bbb/{foo}/{bar}/ccc/ddd/",
			reqPath: "/aaa/bbb/hello/world/ccc/ddd/",
			wantParamMap: map[string]string{
				"foo": "hello",
				"bar": "world",
			},
			wantErr: false,
		},
		{name: "valid_really_touching",
			path:    "/aaa/bbb/{foo}{bar}/ccc/ddd/",
			reqPath: "/aaa/bbb/waaahluigi/ccc/ddd/",
			wantParamMap: map[string]string{
				"foo": "waaahluig",
				"bar": "i",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Start by building the regex
			regex, names, err := buildMatcherRegex(tt.path)
			if err != nil {
				t.Fatalf("Tokenisation integration test: failed to build regex (%s)", err.Error())
				return
			}

			// Now apply that regex to a path
			matches, err := tokenise(regex, names, tt.reqPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenisation Integration Test: tokenise() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantParamMap, matches)

		})
	}
}
