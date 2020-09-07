package aqb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
)

// PlayShare represent play share files
type PlayShare struct {
	GameID       string `json:"game_id"`
	GameName     string `json:"game_name"`
	FilePath     string `json:"file_path"`
	LastPlayDate uint   `json:"last_play_date"`
	ImageWidth   uint   `json:"image_width"`
	ImageHeight  uint   `json:"image_height"`
}

type playShareResponse struct {
	common
	List []*PlayShare `json:"list"`
}

// ListPlayShare lists all playshare
func (c *Client) ListPlayShare() ([]*PlayShare, error) {
	resp, err := c.Get(aqbOrigin + "/aqb/blog/post/webdav/index.php?format=json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var psl playShareResponse
	if err := json.Unmarshal(b, &psl); err != nil {
		return nil, err
	}

	if !psl.IsSuccess {
		return nil, fmt.Errorf("server returns error: %s", psl.Message)
	}

	return psl.List, nil
}

// DownloadPlayShare download PlayShare to writer
func (c *Client) DownloadPlayShare(ps *PlayShare, w io.Writer) error {
	u, err := url.Parse(aqbOrigin + "/aqb/blog/post/webdav/detail.php")
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("filepath", ps.FilePath)
	u.RawQuery = q.Encode()

	res, err := c.Get(u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if _, err := io.Copy(w, res.Body); err != nil {
		return err
	}

	return nil
}
