package ohdear

import "fmt"

type Monitor struct {
	ID     int
	URL    string
	TeamID int `json:"team_id"`
	Checks []Check
}

func (c *Client) GetMonitor(id int) (*Monitor, error) {
	resp, err := c.R().
		SetResult(&Monitor{}).
		Get(fmt.Sprintf("/api/monitors/%d", id))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*Monitor), nil
}

func (c *Client) AddMonitor(url string, teamID int, checks []string) (*Monitor, error) {
	resp, err := c.R().
		SetBody(map[string]interface{}{
			"url":     url,
			"type":    "http",
			"team_id": teamID,
			"checks":  checks,
		}).
		SetResult(&Monitor{}).
		Post("/api/monitors")
	if err != nil {
		return nil, err
	}

	return resp.Result().(*Monitor), nil
}

func (c *Client) RemoveMonitor(id int) error {
	_, err := c.R().Delete(fmt.Sprintf("/api/monitors/%d", id))
	return err
}
