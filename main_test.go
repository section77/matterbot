package main

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/section77/matterbot/logger"
)

// TestMain is the driver for the unit tests
//
// to enable log output:
//    go test -log
//
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

func TestParseFwdMappings(t *testing.T) {
	// single valid value
	validateParseFwdMappings([]fwdMapping{
		fwdMapping{
			marker:   "user",
			mailAddr: "user@mail.com",
		},
	}, "user=user@mail.com", t)

	// pair of valid values
	validateParseFwdMappings([]fwdMapping{
		fwdMapping{"user1", "user1@mail.com"},
		fwdMapping{"user2", "abc@gmail.com"}},
		"user1=user1@mail.com,user2=abc@gmail.com", t)

	// pair with valid values and spaces
	validateParseFwdMappings([]fwdMapping{
		fwdMapping{"user1", "user1@mail.com"},
		fwdMapping{"user2", "abc@gmail.com"}},
		" user1 = user1@mail.com , user2 = abc@gmail.com", t)

	// empty: invalid
	validateParseFwdMappingsErr("empty input", "", "flag 'forward' are mandatory", t)

	// single invalid value
	validateParseFwdMappingsErr("single invalid value", "name",
		"invalid format in flag 'forward': 'name' - valid example: 'user=abc@mail.com'", t)
}

func validateParseFwdMappings(expected []fwdMapping, s string, t *testing.T) {
	res, err := parseFwdMappings(s)
	if err != nil {
		t.Errorf("unexpected error for input: %s", s)
	}

	if len(expected) != len(res) {
		t.Errorf("expected count of elements: %d, received: %d", len(expected), len(res))
	} else {
		for i, e := range expected {
			if e.marker != res[i].marker || e.mailAddr != res[i].mailAddr {
				t.Errorf("didn't match - expected: %+v, received: %+v", e, res[i])
			}
		}
	}
}

func validateParseFwdMappingsErr(name, s, expected string, t *testing.T) {
	_, err := parseFwdMappings(s)
	if err == nil {
		t.Errorf("%s: no error was returned for input: '%s'", name, s)
		return
	}

	if err.Error() != expected {
		t.Errorf("%s: expected error not found - found: \"%s\", expected: \"%s\"", name, err.Error(), expected)
	}
}
