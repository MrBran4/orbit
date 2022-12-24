package orbit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_errRouteDoesNotmatch_Error(t *testing.T) {

	assert.Equal(
		t,
		"route doesn't match (example details)",
		errRouteDoesNotMatch("example details").Error(),
	)

}

func Test_errMisconfigured_Error(t *testing.T) {

	assert.Equal(
		t,
		"orbit may be misconfigured: example details",
		errMisconfigured("example details").Error(),
	)

}

func Test_errCoudlntGetParams_Error(t *testing.T) {

	assert.Equal(
		t,
		"couldn't get paramname from request (example details)",
		errCoudlntGetParams{
			paramName: "paramname",
			err:       errors.New("example details"),
		}.Error(),
	)

}
