package main

import (
	"fmt"
	"os"

	"github.com/otofune/automate-eamusement-playshare/aqb/logger"
)

// debugLogger implements aqb/logger.Logger
type debugLogger struct {
	service string
}

func (d *debugLogger) Debugf(format string, o ...interface{}) (int, error) {
	if d.service != "" {
		format = "[" + d.service + "] " + format
	}
	return os.Stderr.Write([]byte(fmt.Sprintf(format, o...)))
}

func (d debugLogger) WithServiceName(newService string) logger.Logger {
	if d.service != "" {
		d.service += "."
	}
	d.service += newService
	return &d
}
