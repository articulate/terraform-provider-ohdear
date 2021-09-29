package ohdear

import "fmt"

type Check struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func (c *Client) EnableCheck(id int) error {
	_, err := c.R().Post(fmt.Sprintf("/api/checks/%d/enable", id))
	return err
}

func (c *Client) DisableCheck(id int) error {
	_, err := c.R().Post(fmt.Sprintf("/api/checks/%d/disable", id))
	return err
}
