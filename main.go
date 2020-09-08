package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/kelseyhightower/envconfig"
	"github.com/otofune/automate-eamusement-playshare/aqb"
	aqbContext "github.com/otofune/automate-eamusement-playshare/aqb/context"
)

type config struct {
	Username string `required:"true"`
	Password string `required:"true"`
	// following will be given from args (optional, with default)
	StateFile       string `ignored:"true"`
	ImagesDirectory string `ignored:"true"`
}

const (
	defaultStateFile         = "state.bin"
	defaultDownloadDirectory = "images"
)

func main() {
	var c config
	envconfig.MustProcess("AEAP", &c)

	commandName := filepath.Base(os.Args[0])

	parser := argparse.NewParser(commandName, "Download all playshare images from e-AMUSEMENT app server, and then print downloaded path to stdout")
	stateFile := parser.String("s", "state-file", &argparse.Options{
		Help:    fmt.Sprintf("%s will be save state to. Such as session cookie", commandName),
		Default: defaultStateFile,
	})
	imagesDirectory := parser.String("d", "directory", &argparse.Options{
		Help:    "A directory which play share images will be saved",
		Default: defaultDownloadDirectory,
	})
	debugEnabled := parser.Flag("", "debug", &argparse.Options{Help: "Show debug output to stderr"})
	isShowVersion := parser.Flag("v", "version", &argparse.Options{Help: "Show version"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	c.ImagesDirectory = *imagesDirectory
	c.StateFile = *stateFile

	if *isShowVersion {
		if err := showVersion(); err != nil {
			panic(err)
		}
		return
	}

	ctx := context.Background()
	if *debugEnabled {
		ctx = aqbContext.WithLogger(ctx, &debugLogger{})
	}

	if err := run(ctx, c); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, c config) error {
	client, err := aqb.NewClient(ctx)
	if err != nil {
		return err
	}

	if err := restoreOrLogin(ctx, c, &client); err != nil {
		return err
	}

	list, err := client.ListPlayShare()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(c.ImagesDirectory, 0o755); err != nil {
		return err
	}
	for _, item := range list {
		fp := filepath.Join(c.ImagesDirectory, fmt.Sprintf("%s-%d%s", item.GameName, item.LastPlayDate, filepath.Ext(item.FilePath)))
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

	return nil
}
