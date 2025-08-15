package ohdear

import (
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMonitor(t *testing.T) {
	_, reset := mocklog()
	t.Cleanup(reset)
	t.Cleanup(httpmock.DeactivateAndReset)

	resp, err := httpmock.NewJsonResponder(200, map[string]interface{}{
		"id":      1234,
		"url":     "https://example.com",
		"team_id": 5678,
		"checks": []map[string]interface{}{
			{
				"id":      12,
				"type":    "uptime",
				"enabled": true,
			},
			{
				"id":      34,
				"type":    "performance",
				"enabled": false,
			},
		},
	})
	require.NoError(t, err)
	httpmock.RegisterResponder("GET", "https://ohdear.test/api/monitors/1234", resp)

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	monitor, err := client.GetMonitor(1234)
	require.NoError(t, err)
	assert.Equal(t, 1234, monitor.ID)
	assert.Equal(t, "https://example.com", monitor.URL)
	assert.Equal(t, 5678, monitor.TeamID)
	assert.Len(t, monitor.Checks, 2)
	assert.Equal(t, 12, monitor.Checks[0].ID)
	assert.Equal(t, "uptime", monitor.Checks[0].Type)
	assert.True(t, monitor.Checks[0].Enabled)
	assert.Equal(t, "performance", monitor.Checks[1].Type)
	assert.False(t, monitor.Checks[1].Enabled)
}

func TestAddMonitor(t *testing.T) {
	_, reset := mocklog()
	t.Cleanup(reset)
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("POST", "https://ohdear.test/api/monitors",
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			assert.JSONEq(
				t,
				`{"checks":["uptime","performance","broken_links"],"team_id":5678,"url":"https://example.com/new","type":"http"}`,
				string(body),
			)

			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"id":      4321,
				"url":     "https://example.com/new",
				"team_id": 5678,
				"checks": []map[string]interface{}{
					{
						"id":      12,
						"type":    "uptime",
						"enabled": true,
					},
				},
			})
		})

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	monitor, err := client.AddMonitor("https://example.com/new", 5678, []string{"uptime", "performance", "broken_links"})
	require.NoError(t, err)
	assert.Equal(t, 4321, monitor.ID)
	assert.Equal(t, "https://example.com/new", monitor.URL)
}

func TestRemoveMonitor(t *testing.T) {
	_, reset := mocklog()
	t.Cleanup(reset)
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("DELETE", "https://ohdear.test/api/monitors/1234", httpmock.NewStringResponder(204, ""))

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	err := client.RemoveMonitor(1234)
	require.NoError(t, err)
}
