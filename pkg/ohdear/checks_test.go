package ohdear

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestEnableCheck(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

	resp, err := httpmock.NewJsonResponder(200, map[string]interface{}{"id": 1234})
	assert.NoError(t, err)
	httpmock.RegisterResponder("POST", "https://ohdear.test/api/checks/1234/enable", resp)

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	err = client.EnableCheck(1234)
	assert.NoError(t, err)
}

func TestDisableCheck(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

	resp, err := httpmock.NewJsonResponder(200, map[string]interface{}{"id": 4321})
	assert.NoError(t, err)
	httpmock.RegisterResponder("POST", "https://ohdear.test/api/checks/4321/disable", resp)

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	err = client.DisableCheck(4321)
	assert.NoError(t, err)
}
