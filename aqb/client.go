package aqb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	aqbOrigin            = "https://aqb.s.konaminet.jp"
	persistentCookieName = "aqblog"
	eamLinkIOSVersion    = "3.5.2.59"
	apiUserAgent         = "jp.konami.eam.link (iPhone12,8; iOS 13.6.1; in-app; 20; app-version; " + eamLinkIOSVersion + ")"
	storageUserAgent     = "eAMUSEMENT/59 CFNetwork/1128.0.1 Darwin/19.6.0"
)

// Client interact with aqb server
type Client struct {
	*http.Client
}

type common struct {
	// false when error
	IsSuccess  bool `json:"status"`
	StatusCode uint `json:"status_code"`
	// message may be included when status = false
	Message string
}

type uaRoundTripper struct {
	rt http.RoundTripper
}

func (ur *uaRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", apiUserAgent)

	if r.URL.Path == "/aqb/blog/post/webdav/detail.php" {
		r.Header.Set("User-Agent", storageUserAgent)
	}

	fmt.Printf("[request] %s\n", r.URL)

	return ur.rt.RoundTrip(r)
}

// NewClient build new Client
func NewClient() (Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return Client{}, err
	}

	t := &uaRoundTripper{rt: http.DefaultTransport}
	return Client{Client: &http.Client{Jar: jar, Transport: t}}, nil
}

func (c *Client) postForm(path string, form url.Values) (*[]byte, error) {
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", aqbOrigin+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[postForm] %s(%d): %s\n", path, res.StatusCode, buf)
	fmt.Printf("\treq header: %#v\n", res.Request.Header)
	fmt.Printf("\tres header: %#v\n", res.Header)

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Server returns error status code(%d): %s", res.StatusCode, buf)
	}

	return &buf, nil
}
