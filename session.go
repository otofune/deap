package main

import (
	"context"
	"io"
	"os"

	"github.com/otofune/automate-eamusement-playshare/aqb"
	aqbContext "github.com/otofune/automate-eamusement-playshare/aqb/context"
)

func loginAndSaveSession(ctx context.Context, conf config, client *aqb.Client) error {
	if err := client.Login(conf.Username, conf.Password); err != nil {
		return err
	}
	r, err := client.GetSession()
	if err != nil {
		return err
	}
	w, err := os.OpenFile(conf.StateFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

func restoreOrLogin(ctx context.Context, conf config, client *aqb.Client) error {
	logger := aqbContext.Logger(ctx).WithServiceName("restoreOrLogin")

	f, err := os.Open(conf.StateFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// try to login
	if err == nil {
		logger.Debugf("restores state from state file\n")
		defer f.Close()

		if err := client.RestoreSession(f); err != nil {
			return err
		}

		ok, err := client.CheckSession()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		logger.Debugf("checkSession failed, try to login\n")
	}

	// try to login
	if err := loginAndSaveSession(ctx, conf, client); err != nil {
		return err
	}

	return nil
}
