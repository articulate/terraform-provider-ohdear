package ohdear

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Run("error with message and errors", func(t *testing.T) {
		err := &Error{
			Message: "test error",
			Errors: map[string][]string{
				"foo": {"bar", "baz"},
			},
		}

		assert.EqualError(t, err, "test error\nfoo: bar baz")
	})

	t.Run("error with no errors", func(t *testing.T) {
		err := &Error{
			Message: "test error",
		}

		assert.EqualError(t, err, "test error")
	})

	t.Run("error with no message", func(t *testing.T) {
		err := &Error{
			Response: &resty.Response{
				RawResponse: &http.Response{
					Status:     "Unauthorized",
					StatusCode: 401,
				},
			},
		}

		assert.EqualError(t, err, "401: Unauthorized")
	})
}

func TestErrorFromResponse(t *testing.T) {
	_, reset := mocklog()
	defer reset()
	defer httpmock.DeactivateAndReset()

	resp, err := httpmock.NewJsonResponder(404, map[string]interface{}{"message": "Not found"})
	assert.NoError(t, err)
	httpmock.RegisterResponder("GET", "https://ohdear.test/api/sites/1", resp)

	client := NewClient("https://ohdear.test", "")
	httpmock.ActivateNonDefault(client.GetClient())

	_, err = client.R().Get("/api/sites/1")
	assert.Error(t, err)

	var e *Error
	assert.True(t, errors.As(err, &e))
	assert.Equal(t, 404, e.Response.StatusCode())
}
