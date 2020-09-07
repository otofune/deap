package aqb

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	aqbOrigin            = "https://aqb.s.konaminet.jp"
	persistentCookieName = "aqblog"
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

// NewClient establish
func NewClient() (Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return Client{}, err
	}
	return Client{Client: &http.Client{Jar: jar}}, nil
}

// GetSession returns io.Reader that includes information to recovery session.
// Internally returns value of aqblog= cookie currently.
// That cookie is the only persistent cookie returned from server.
func (c *Client) GetSession() (io.Reader, error) {
	u, err := url.Parse(aqbOrigin)
	if err != nil {
		return strings.NewReader(""), fmt.Errorf("unexpected error happened, must be unreachable: %w", err)
	}

	cookies := c.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == persistentCookieName {
			return strings.NewReader(cookie.Value), nil
		}
	}

	fmt.Println("[GetSession] No " + persistentCookieName)
	return strings.NewReader(""), nil
}

// RestoreSession recovery session from GetSession return value.
func (c *Client) RestoreSession(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	str := string(b)
	if str == "" {
		return nil
	}

	u, err := url.Parse(aqbOrigin)
	if err != nil {
		return fmt.Errorf("unexpected error happened, must be unreachable: %w", err)
	}

	cookie := http.Cookie{
		Name:  persistentCookieName,
		Value: str,
		// same as server returns
		MaxAge:   60 * 60 * 24 * 300,
		Domain:   u.Host,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	c.Jar.SetCookies(u, []*http.Cookie{&cookie})

	return nil
}
