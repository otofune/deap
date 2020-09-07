package main

import (
	"fmt"
	"io"
	"os"

	"github.com/otofune/automate-eamusement-playshare/aqb"
)

func loginAndSaveSession(conf config, client *aqb.Client) error {
	if err := client.Login(conf.Username, conf.Password); err != nil {
		return err
	}
	r, err := client.GetSession()
	if err != nil {
		return err
	}
	w, err := os.OpenFile(conf.SessionFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func restoreOrLogin(conf config, client *aqb.Client) error {
	f, err := os.Open(conf.SessionFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// try to restore
	if err == nil {
		fmt.Println("[restoreOrLogin] restores")
		defer f.Close()
		return client.RestoreSession(f)
	}

	// try to login
	if err := loginAndSaveSession(conf, client); err != nil {
		return err
	}

	return nil
}
