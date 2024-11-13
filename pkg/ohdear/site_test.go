package ohdear

import (
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSite(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

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
	httpmock.RegisterResponder("GET", "https://ohdear.test/api/sites/1234", resp)

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	site, err := client.GetSite(1234)
	require.NoError(t, err)
	assert.Equal(t, 1234, site.ID)
	assert.Equal(t, "https://example.com", site.URL)
	assert.Equal(t, 5678, site.TeamID)
	assert.Len(t, site.Checks, 2)
	assert.Equal(t, 12, site.Checks[0].ID)
	assert.Equal(t, "uptime", site.Checks[0].Type)
	assert.True(t, site.Checks[0].Enabled)
	assert.Equal(t, "performance", site.Checks[1].Type)
	assert.False(t, site.Checks[1].Enabled)
}

func TestAddSite(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://ohdear.test/api/sites",
		func(req *http.Request) (*http.Response, error) {
			body, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			assert.JSONEq(
				t,
				`{"checks":["uptime","performance","broken_links"],"team_id":5678,"url":"https://example.com/new"}`,
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

	site, err := client.AddSite("https://example.com/new", 5678, []string{"uptime", "performance", "broken_links"})
	require.NoError(t, err)
	assert.Equal(t, 4321, site.ID)
	assert.Equal(t, "https://example.com/new", site.URL)
}

func TestRemoveSite(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("DELETE", "https://ohdear.test/api/sites/1234", httpmock.NewStringResponder(204, ""))

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	err := client.RemoveSite(1234)
	require.NoError(t, err)
}
