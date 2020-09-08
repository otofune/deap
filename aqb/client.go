package aqb

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/otofune/deap/aqb/aqbctx"
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
	ctx context.Context
}

type common struct {
	// false when error
	IsSuccess  bool `json:"status"`
	StatusCode uint `json:"status_code"`
	// message may be included when status = false
	Message string
}

type uaRoundTripper struct {
	rt  http.RoundTripper
	ctx context.Context
}

func (ur *uaRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	logger := aqbctx.Logger(ur.ctx).WithServiceName("request")

	r.Header.Set("User-Agent", apiUserAgent)

	if r.URL.Path == "/aqb/blog/post/webdav/detail.php" {
		r.Header.Set("User-Agent", storageUserAgent)
	}

	logger.Debugf("%s %s\n", r.Method, r.URL)

	return ur.rt.RoundTrip(r)
}

// NewClient build new Client
func NewClient(ctx context.Context) (Client, error) {
	ctx = aqbctx.WithLogger(ctx, aqbctx.Logger(ctx).WithServiceName("aqb"))

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return Client{}, err
	}

	t := &uaRoundTripper{rt: http.DefaultTransport, ctx: ctx}
	return Client{ctx: ctx, Client: &http.Client{Jar: jar, Transport: t}}, nil
}

func (c *Client) postForm(path string, form url.Values) (*[]byte, error) {
	logger := aqbctx.Logger(c.ctx).WithServiceName("postForm")

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

	logger.Debugf("%s(%d): %s\n", path, res.StatusCode, buf)
	logger.Debugf("\treq header: %#v\n", res.Request.Header)
	logger.Debugf("\tres header: %#v\n", res.Header)

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Server returns error status code(%d): %s", res.StatusCode, buf)
	}

	return &buf, nil
}
