package orbit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testBodyableTypeStruct struct {
	FieldOne string `json:"field_one"`
	FieldTwo int    `json:"field_two"`
}

func (ts testBodyableTypeStruct) FromBody(body io.ReadCloser) (any, error) {
	var result testBodyableTypeStruct
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("couldn't decode json")
	}
	return result, nil
}

type testBodyableTypeStructReturnsWrongType struct {
	FieldOne string `json:"field_one"`
	FieldTwo int    `json:"field_two"`
}

// This type's FromBody incorrectly returns a testBodyableTypeStruct
// (instead of a testBodyableTypeStructReturnsWrongType)
func (ts testBodyableTypeStructReturnsWrongType) FromBody(body io.ReadCloser) (any, error) {
	var result testBodyableTypeStruct
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("couldn't decode json")
	}
	return result, nil
}

func Test_FromBodyable_Valid(t *testing.T) {

	inputBuf := io.NopCloser(bytes.NewBufferString(`{
		"field_one": "Hello World",
		"field_two": 128
	}`))

	expectedType := testBodyableTypeStruct{}
	expected := testBodyableTypeStruct{
		FieldOne: "Hello World",
		FieldTwo: 128,
	}

	result, err := tryFromBody(expectedType, inputBuf)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

}

func Test_FromBodyable_NotValid(t *testing.T) {

	inputBuf := io.NopCloser(bytes.NewBufferString(`{
		"field_o Not valid json
	}`))

	expectedType := testBodyableTypeStruct{}

	_, err := tryFromBody(expectedType, inputBuf)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "couldn't decode json")

}

func Test_FromBodyable_MismatchedType(t *testing.T) {

	// this'll successfully decode, but its FromBody handler returns it as
	// the wrong type, which should cause the handler to halt.
	inputBuf := io.NopCloser(bytes.NewBufferString(`{
		"field_one": "Hello World",
		"field_two": 128
	}`))

	expectedType := testBodyableTypeStructReturnsWrongType{}

	_, err := tryFromBody(expectedType, inputBuf)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "unexpected type")

}
