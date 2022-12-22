package orbit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newFromRequest_Valid(t *testing.T) {

	// setup
	var paramTypes = RouteParams{
		"stringparam": testTypeString(""),
		"intparam":    testTypeInt(0),
		"structparam": testTypeStruct{},
	}

	var expected = RouteParams{
		"stringparam": testTypeString("hello"),
		"intparam":    testTypeInt(12345),
		"structparam": testTypeStruct{valuePassedIn: "world"},
	}

	// do
	result, err := paramTypes.newFromRequest(map[string]string{
		"stringparam": "hello",
		"intparam":    "12345",
		"structparam": "world",
	})

	// check
	assert.NoError(t, err)
	assert.Equal(t, expected, *result)

}

func Test_newFromRequest_Unparsable(t *testing.T) {

	// setup
	var paramTypes = RouteParams{
		"stringparam": testTypeString(""),
		"intparam":    testTypeInt(0),
		"structparam": testTypeStruct{},
	}

	// do
	_, err := paramTypes.newFromRequest(map[string]string{
		"stringparam": "hello",
		"intparam":    "NOT_AN_INT",
		"structparam": "world",
	})

	// check
	assert.Error(t, err)

}

func Test_newFromRequest_WrongLength(t *testing.T) {

	// setup
	var paramTypes = RouteParams{
		"stringparam": testTypeString(""),
		"intparam":    testTypeInt(0),
		"structparam": testTypeStruct{},
	}

	// do
	_, err := paramTypes.newFromRequest(map[string]string{
		"stringparam": "hello",
		"structparam": "world",
	})

	// check
	assert.Error(t, err)

}

func Test_newFromRequest_BadTypes(t *testing.T) {

	// setup
	var paramTypes = RouteParams{
		"stringparam": testTypeString(""),
		"intparam":    testTypeBadIntWillFail(0),
		"structparam": testTypeStruct{},
	}

	// do
	_, err := paramTypes.newFromRequest(map[string]string{
		"stringparam": "hello",
		"intparam":    "12345",
		"structparam": "world",
	})

	// check
	assert.Error(t, err)

}
