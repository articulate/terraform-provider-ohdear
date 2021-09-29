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
		Get(fmt.Sprintf("/api/sites/%d", id))

	if err != nil {
		return nil, err
	}

	return resp.Result().(*Site), nil
}

func (c *Client) AddSite(url string, teamID int, checks []string) (*Site, error) {
	resp, err := c.R().
		SetBody(map[string]interface{}{
			"url":     url,
			"team_id": teamID,
			"checks":  checks,
		}).
		SetResult(&Site{}).
		Post("/api/sites")

	if err != nil {
		return nil, err
	}

	return resp.Result().(*Site), nil
}

func (c *Client) RemoveSite(id int) error {
	_, err := c.R().Delete(fmt.Sprintf("/api/sites/%d", id))
	return err
}
