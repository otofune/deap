package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	c.SessionFile = "state.bin" // TODO: accept from arg or env

	client, err := aqb.NewClient()
	if err != nil {
		panic(err)
	}

	if err := restoreOrLogin(c, &client); err != nil {
		panic(err)
	}

	list, err := client.ListPlayShare()
	if err != nil {
		panic(err)
	}

	for _, item := range list {
		// TODO: accept images/ from arg or env
		fp := "images/" + fmt.Sprintf("%s-%d%s", item.GameName, item.LastPlayDate, filepath.Ext(item.FilePath))
		if _, err := os.Open(fp); !os.IsNotExist(err) {
			if err != nil {
				fmt.Printf("%s\n", err)
			}
			continue
		}

		file, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		defer file.Close()

		if err := client.DownloadPlayShare(item, file); err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
