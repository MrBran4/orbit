package orbit

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Dummy string type implementing FromRequest
type testTypeString string

func (x *testTypeString) FromRequest(param string) error {
	val := testTypeString(param)
	*x = val
	return nil
}

// Dummy int32 type implementing FromRequest
type testTypeInt int32

func (x *testTypeInt) FromRequest(param string) error {
	intval, err := strconv.Atoi(param)
	if err != nil {
		return err
	}

	result := testTypeInt(intval)
	*x = result
	return nil
}

// Dummy complex type implementing FromRequest
type testTypeStruct struct {
	valuePassedIn string
}

func (x *testTypeStruct) FromRequest(param string) error {
	result := testTypeStruct{
		valuePassedIn: param,
	}
	*x = result
	return nil
}

func Test_FromRequest_testTypeString(t *testing.T) {

	// setup
	var result testTypeString

	// do
	err := result.FromRequest("hello there")

	// check
	assert.NoError(t, err)
	assert.Equal(t, testTypeString("hello there"), result)

}

func Test_FromRequest_testTypeInt(t *testing.T) {

	// setup
	var result testTypeInt

	// do
	err := result.FromRequest("128")

	// check
	assert.NoError(t, err)
	assert.Equal(t, testTypeInt(128), result)

}

func Test_FromRequest_testTypeStruct(t *testing.T) {

	// setup
	var result testTypeStruct

	// do
	err := result.FromRequest("oof")

	// check
	assert.NoError(t, err)
	assert.Equal(t, testTypeStruct{valuePassedIn: "oof"}, result)

}
