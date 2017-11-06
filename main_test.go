package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/section77/matterbot/logger"
)

func TestMain(m *testing.M) {
	logOutput := flag.Bool("log", false, "enable logging to stdout")
	flag.Parse()

	if *logOutput {
		logger.SetLogLevel(logger.DebugLevel)
	} else {
		logger.SetLogLevel(logger.Disabled)
	}
	//os.Exit(m.Run())
	res := m.Run()
	time.Sleep(500 * time.Millisecond)
	os.Exit(res)

}
