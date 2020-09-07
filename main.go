package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/otofune/automate-eamusement-playshare/aqb"
)

type config struct {
	Username    string `required:"true"`
	Password    string `required:"true"`
	SessionFile string `ignored:"true"`
}

func main() {
	var c config
	envconfig.MustProcess("AEAP", &c)
	c.SessionFile = "state.bin" // TODO: accept change from arg or env

	client, err := aqb.NewClient()
	if err != nil {
		panic(err)
	}

	if err := restoreOrLogin(c, &client); err != nil {
		panic(err)
	}
}
