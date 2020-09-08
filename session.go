package main

import (
	"context"

	"github.com/otofune/automate-eamusement-playshare/aqb"
	"github.com/otofune/automate-eamusement-playshare/aqb/aqbctx"
)

func checkSessionOrLogin(ctx context.Context, client *aqb.Client, username, password string) error {
	logger := aqbctx.Logger(ctx)

	ok, err := client.CheckSession()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	logger.Debugf("try to login\n")

	if err := client.Login(username, password); err != nil {
		return err
	}

	return nil
}
