package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/akamensky/argparse"
	"github.com/kelseyhightower/envconfig"
	"github.com/otofune/deap/aqb"
	"github.com/otofune/deap/aqb/aqbctx"
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

// main load environment variables & arguments to config and run command with config.
// If required, main loads & save state too.
func main() {
	var c config
	envconfig.MustProcess("DEAP", &c)

	commandName := filepath.Base(os.Args[0])

	parser := argparse.NewParser(commandName, "Download all playshare images from e-AMUSEMENT app server, and then print downloaded path to stdout")
	stateFile := parser.String("s", "state-file", &argparse.Options{
		Help:    fmt.Sprintf("%s will be save state to, such as session cookie", commandName),
		Default: defaultStateFile,
	})
	imagesDirectory := parser.String("d", "directory", &argparse.Options{
		Help:    fmt.Sprintf("A directory which play share images will be saved. %s will download not-downloaded images", commandName),
		Default: defaultDownloadDirectory,
	})
	debugEnabled := parser.Flag("", "debug", &argparse.Options{Help: "Show debug output to stderr"})
	isShowVersion := parser.Flag("v", "version", &argparse.Options{Help: "Show version"})

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	c.ImagesDirectory = *imagesDirectory
	c.StateFile = *stateFile

	if *isShowVersion {
		if err := showVersion(); err != nil {
			panic(err)
		}
		return
	}

	logger := debugLogger{}

	ctx := context.Background()
	if *debugEnabled {
		ctx = aqbctx.WithLogger(ctx, &logger)
	}

	s := state{}
	if err := s.loadState(ctx, c); err != nil {
		logger.Errorf("%s\n", err)
		os.Exit(1)
	}

	statusCode := 0
	if err := run(ctx, c, &s); err != nil {
		logger.Errorf("%s\n", err)
		statusCode = 1
	}
	if err := s.saveState(ctx, c); err != nil {
		logger.Errorf("%s\n", err)
		statusCode = 1
	}
	os.Exit(statusCode)
}

func run(ctx context.Context, c config, s *state) error {
	client, err := aqb.NewClient(ctx)
	if err != nil {
		return err
	}

	if s.AQBSession != nil {
		if err := client.RestoreSession(bytes.NewReader(*s.AQBSession)); err != nil {
			return err
		}
	}
	if err := checkSessionOrLogin(ctx, &client, c.Username, c.Password); err != nil {
		return err
	}
	// save to state
	r, err := client.GetSession()
	if err != nil {
		return err
	}
	rb, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.AQBSession = (*base64bytes)(&rb)

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
				return err
			}
			// skip if existed
			continue
		}

		file, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer file.Close()

		if err := client.DownloadPlayShare(item, file); err != nil {
			return err
		}
		fmt.Println(fp)
	}

	return nil
}
