package ohdear

import "fmt"

type Site struct {
	ID     int
	URL    string
	TeamID int `json:"team_id"`
	Checks []Check
}

func (c *Client) GetSite(id int) (*Site, error) {
	resp, err := c.R().
		SetResult(&Site{}).
		Get(fmt.Sprintf("/api/monitors/%d", id))
	if err != nil {
		return nil, fmt.Errorf("could not get site %d: %w", id, err)
	}

	return resp.Result().(*Site), nil
}

func (c *Client) AddSite(url string, teamID int, checks []string) (*Site, error) {
	resp, err := c.R().
		SetBody(map[string]any{
			"url":     url,
			"type":    "http",
			"team_id": teamID,
			"checks":  checks,
		}).
		SetResult(&Site{}).
		Post("/api/monitors")
	if err != nil {
		return nil, fmt.Errorf("could not add site: %w", err)
	}

	return resp.Result().(*Site), nil
}

func (c *Client) RemoveSite(id int) error {
	if _, err := c.R().Delete(fmt.Sprintf("/api/monitors/%d", id)); err != nil {
		return fmt.Errorf("could not remove site %d: %w", id, err)
	}
	return nil
}
