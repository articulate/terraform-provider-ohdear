package ohdear

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Error struct {
	Response *resty.Response `json:"-"`
	Message  string
	Errors   map[string][]string
}

func (e Error) Error() string {
	if e.Message != "" {
		msg := e.Message
		for key, err := range e.Errors {
			msg = fmt.Sprintf("%s\n%s: %s", msg, key, strings.Join(err, " "))
		}
		return msg
	}

	return fmt.Sprintf("%d: %s", e.Response.StatusCode(), e.Response.Status())
}

func errorFromResponse(_ *resty.Client, r *resty.Response) error {
	if !r.IsError() {
		return nil
	}

	err := r.Error().(*Error)
	err.Response = r
	return err
}
