package orbit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicString_FromRequest(t *testing.T) {

	var paramType = BasicString("")

	result, err := paramType.FromRequest("example_string")
	assert.NoError(t, err)

	decoded, _ := result.(BasicString)
	assert.Equal(t, "example_string", string(decoded))

}

func Test_BasicInt_FromRequest_Valid(t *testing.T) {

	var paramType = BasicInt(0)

	result, err := paramType.FromRequest("12345")
	assert.NoError(t, err)

	decoded, _ := result.(BasicInt)
	assert.Equal(t, 12345, int(decoded))

}

func Test_BasicInt_FromRequest_Invalid(t *testing.T) {

	var paramType = BasicInt(0)

	_, err := paramType.FromRequest("aaaaa")
	assert.Error(t, err)

}
