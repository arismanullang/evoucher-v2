// +build linux,syslog

package main

import (
	"log"
	"log/syslog"
)

func init() {
	logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_DAEMON, name)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logger)
	log.SetFlags(0)
}
