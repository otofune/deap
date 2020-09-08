package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
)

// base64bytes is []byte in Go, but string in JSON
type base64bytes []byte

func (bb base64bytes) MarshalJSON() ([]byte, error) {
	s := base64.StdEncoding.EncodeToString(bb)
	return json.Marshal(s)
}

func (bb *base64bytes) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	var err error
	*bb, err = base64.StdEncoding.DecodeString(s)
	return err
}

type state struct {
	AQBSession *base64bytes
}

// loadState loads state from stateFile if stateFile exists
func (s *state) loadState(ctx context.Context, conf config) error {
	f, err := os.Open(conf.StateFile)
	if err != nil {
		// new
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, s); err != nil {
		return err
	}

	return nil
}

func (s *state) saveState(ctx context.Context, conf config) error {
	f, err := os.OpenFile(conf.StateFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}
