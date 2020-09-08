package main

import (
	"fmt"
	"log"
	"os"

	"github.com/otofune/automate-eamusement-playshare/aqb/logger"
)

// debugLogger implements aqb/logger.Logger
type debugLogger struct {
	service      string
	stderrLogger *log.Logger
}

func (d *debugLogger) getStderrLogger() *log.Logger {
	if d.stderrLogger == nil {
		d.stderrLogger = log.New(os.Stderr, "", log.Lshortfile|log.Ldate)
	}
	return d.stderrLogger
}

func (d *debugLogger) Debugf(format string, o ...interface{}) {
	d.getStderrLogger().Output(2, fmt.Sprintf(format, o...))
}

func (d *debugLogger) Errorf(format string, o ...interface{}) {
	d.getStderrLogger().Output(2, fmt.Sprintf(format, o...))
}

func (d debugLogger) WithServiceName(newService string) logger.Logger {
	if d.service != "" {
		d.service += "."
	}
	d.service += newService
	return &d
}
