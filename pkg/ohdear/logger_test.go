package ohdear

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mocklog() (*bytes.Buffer, func()) {
	original := log.Writer()

	output := &bytes.Buffer{}
	log.SetOutput(output)

	return output, func() {
		log.SetOutput(original)
	}
}

func TestTerraformLogger(t *testing.T) {
	logger := &TerraformLogger{}
	out, reset := mocklog()
	t.Cleanup(reset)

	logger.Errorf("test error message")
	assert.Contains(t, out.String(), "[ERROR] test error message\n", out.String())

	logger.Warnf("test warn message")
	assert.Contains(t, out.String(), "[WARN] test warn message\n", out.String())

	logger.Debugf("test debug message")
	assert.Contains(t, out.String(), "[DEBUG] test debug message\n", out.String())
}
