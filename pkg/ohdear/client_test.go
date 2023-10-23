package ohdear

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	token := os.Getenv("OHDEAR_TOKEN")
	teamID := os.Getenv("OHDEAR_TEAM_ID")
	if token == "" || teamID == "" || os.Getenv("TF_ACC") == "" {
		t.Skip("Integration tests skipped unless 'OHDEAR_TOKEN', 'OHDEAR_TEAM_ID', and 'TF_ACC' env vars are set")
	}

	if os.Getenv("SKIP_INTEGRATION_TESTS") != "" {
		t.Skip("Integration tests skipped because SKIP_INTEGRATION_TESTS env var set")
	}

	team, err := strconv.Atoi(teamID)
	require.NoError(t, err)

	url := os.Getenv("OHDEAR_API_URL")
	if url == "" {
		url = "https://ohdear.app"
	}

	_, reset := mocklog()
	// if we defer, we get log leakage from the other cleanup function
	t.Cleanup(reset)

	ua := "terraform-provider-ohdear/TEST (https://github.com/articulate/terraform-provider-ohdear) integration-tests"
	client := NewClient(url, token)
	client.SetDebug(false)
	client.SetUserAgent(ua)

	create, err := client.AddSite("https://example.com", team, []string{"uptime"})
	require.NoError(t, err)

	// make sure we remove the site even if tests fail
	t.Cleanup(func() {
		var e *Error
		if !errors.As(err, &e) {
			t.Fatal("site was not removed from OhDear")
		}

		if e.Response.StatusCode() != 404 {
			t.Fatal("site was not removed from Oh Dear")
		}
	})

	uptime, enabled := getCheckInfo(create)

	assert.Equal(t, "https://example.com", create.URL)
	assert.ElementsMatch(t, []string{"uptime"}, enabled)

	// get the site
	site, err := client.GetSite(create.ID)
	require.NoError(t, err)
	assert.Equal(t, create, site)

	// disable the uptime check
	err = client.DisableCheck(uptime)
	require.NoError(t, err)
	update, err := client.GetSite(site.ID)
	require.NoError(t, err)
	_, enabled = getCheckInfo(update)
	assert.Empty(t, enabled)

	// enable the uptime check
	err = client.EnableCheck(uptime)
	require.NoError(t, err)
	update, err = client.GetSite(site.ID)
	require.NoError(t, err)
	_, enabled = getCheckInfo(update)
	assert.ElementsMatch(t, []string{"uptime"}, enabled)

	// delete the site
	err = client.RemoveSite(site.ID)
	require.NoError(t, err)

	// verify it was deleted (wait because sometimes it takes the api a second to update)
	time.Sleep(5 * time.Second)
	removed, err := client.GetSite(site.ID)
	assert.Nil(t, removed)

	var e *Error
	require.ErrorAs(t, err, &e)
	assert.Equal(t, 404, e.Response.StatusCode())
}

func TestSetUserAgent(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://ohdear.test/ping",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "application/json", req.Header.Get("Accept"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, "user-agent/test", req.Header.Get("User-Agent"))

			return httpmock.NewStringResponse(200, ""), nil
		})

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	client.SetUserAgent("user-agent/test")
	_, err := client.R().Get("/ping")
	require.NoError(t, err)
}

func getCheckInfo(s *Site) (int, []string) {
	uptime := 0
	enabled := []string{}
	for _, check := range s.Checks {
		if check.Enabled {
			enabled = append(enabled, check.Type)
		}
		if check.Type == "uptime" {
			uptime = check.ID
		}
	}

	return uptime, enabled
}
