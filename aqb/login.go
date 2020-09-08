package aqb

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Login issues session
func (c *Client) Login(username, password string) error {
	form := url.Values{}
	form.Add("format", "json")
	form.Add("username", username)
	form.Add("password", password)

	buf, err := c.postForm("/aqb/user/login.php", form)
	if err != nil {
		return err
	}

	var body common
	if err := json.Unmarshal(*buf, &body); err != nil {
		return err
	}

	if !body.IsSuccess {
		return fmt.Errorf("server returns error: %s", body.Message)
	}

	return nil
}

// CheckSession returns
func (c *Client) CheckSession() (bool, error) {
	form := url.Values{}
	form.Add("format", "json")

	buf, err := c.postForm("/aqb/user/checkSession.php", form)
	if err != nil {
		return false, err
	}

	var body common
	if err := json.Unmarshal(*buf, &body); err != nil {
		return false, err
	}

	return body.IsSuccess, nil
}
