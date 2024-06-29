// This is a placeholder for real testing:  the CI workflow fails when the
// coverage output file is empty, and, other than this, this repo currently
// has -no- tests...so this placeholder allows the CI to pass.  It also ensures
// that things at least compile.

package main_test

import (
	"testing"

	service "arcaflow-plugin-service"

	"go.arcalot.io/assert"
)

// TestServiceSchema is a trivial test which makes sure that the Service Schema
// contains the expected step -- mostly, this just ensures that the Service
// sources actually compile.
func TestServiceSchema(t *testing.T) {
	s := service.Schema
	assert.MapContainsKey(t, "create", s.StepsValue)
}
