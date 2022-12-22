package orbit

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Dummy string type implementing FromRequest
type testTypeString string

func (x testTypeString) FromRequest(param string) (any, error) {
	val := testTypeString(param)
	return val, nil
}

// Dummy int32 type implementing FromRequest
type testTypeInt int32

func (x testTypeInt) FromRequest(param string) (any, error) {
	intval, err := strconv.Atoi(param)
	if err != nil {
		return nil, err
	}

	result := testTypeInt(intval)
	return result, nil
}

// Dummy int32 type implementing FromRequest
type testTypeBadIntWillFail int32

func (x testTypeBadIntWillFail) FromRequest(param string) (any, error) {
	result := testTypeString(param)
	return result, nil
}

// Dummy complex type implementing FromRequest
type testTypeStruct struct {
	valuePassedIn string
}

func (x testTypeStruct) FromRequest(param string) (any, error) {
	result := testTypeStruct{
		valuePassedIn: param,
	}
	return result, nil
}

func Test_FromRequest_testTypeString(t *testing.T) {

	// setup
	var zeroedString testTypeString

	// do
	result, err := zeroedString.FromRequest("hello there")

	// check
	assert.NoError(t, err)
	assert.Equal(t, testTypeString("hello there"), result)

}

func Test_FromRequest_testTypeInt(t *testing.T) {

	// setup
	var zeroedInt testTypeInt

	// do
	result, err := zeroedInt.FromRequest("128")

	// check
	assert.NoError(t, err)
	assert.Equal(t, testTypeInt(128), result)

}

func Test_FromRequest_testTypeStruct(t *testing.T) {

	// setup
	var zeroedStruct testTypeStruct

	// do
	result, err := zeroedStruct.FromRequest("oof")

	// check
	assert.NoError(t, err)
	assert.Equal(t, testTypeStruct{valuePassedIn: "oof"}, result)

}
