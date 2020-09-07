package aqb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const urlEncodedType = "application/x-www-form-urlencoded"

// Login issues session
func (c *Client) Login(username, password string) error {
	form := url.Values{}
	form.Add("format", "json")
	form.Add("username", username)
	form.Add("password", password)

	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", aqbOrigin+"/aqb/user/login.php", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBuf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", bodyBuf)
	var bdy common
	if err := json.Unmarshal(bodyBuf, &bdy); err != nil {
		return err
	}

	if !bdy.IsSuccess {
		return fmt.Errorf("server returns error: %s", bdy.Message)
	}

	return nil
}

// CheckSession returns
func (c *Client) CheckSession() (bool, error) {
	form := url.Values{}
	form.Add("format", "json")
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", aqbOrigin+"/aqb/user/checkSession.php", body)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	bodyBuf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	fmt.Printf("%s\n", bodyBuf)
	var bdy common
	if err := json.Unmarshal(bodyBuf, &bdy); err != nil {
		return false, err
	}

	return bdy.IsSuccess, nil
}
