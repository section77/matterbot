package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/namsral/flag"

	"github.com/section77/matterbot/chat"
	"github.com/section77/matterbot/logger"
	"github.com/section77/matterbot/mail"
)

// flags
var (
	logVerbose  = flag.Bool("v", false, "enable verbose / debug output")
	logDisabled = flag.Bool("q", false, "disable logging / be quite")

	showVersion = flag.Bool("V", false, "show version and exit")

	mattermostURL  = flag.String("mattermost-url", "http://127.0.0.1:8065", "mattermost url (http://x.x.x.x:xxxx)")
	mattermostUser = flag.String("mattermost-user", "matterbot", "mattermost user")
	mattermostPass = flag.String("mattermost-pass", "tobrettam", "mattermost password")

	mailHost   = flag.String("mail-host", "127.0.0.1:25", "mail-server host (<HOST>:<PORT>)")
	mailUser   = flag.String("mail-user", "matterbot@localhost", "mail login user")
	mailPass   = flag.String("mail-pass", "tobrettam", "mail login pass")
	mailUseTLS = flag.Bool("mail-use-tls", false, "use TLS instead of STARTTLS")

	mailSubject = flag.String("mail-subject",
		"Forwarding from mattermost: {{.User}} writes in channel {{.Channel}}",
		"mail subject")
	mailBody = flag.String("mail-body",
		"{{.Content}}",
		"mail body")

	forward = flag.String("forward", "",
		"mapping from marker to receiver mail address. example: 'user1=user1@gmail.com,user2=abc@mail.com'")
)

var mailSubjectTemplate *template.Template
var mailBodyTemplate *template.Template

var version string

func main() {

	//
	// parse and process commad line args
	//
	flag.Parse()

	if *showVersion {
		println("matterbot: v" + version)
		os.Exit(0)
	}

	if *logVerbose {
		logger.SetLogLevel(logger.DebugLevel)
	} else if *logDisabled {
		logger.SetLogLevel(logger.ErrorLevel)
	}

	url, err := url.Parse(*mattermostURL)
	if err != nil {
		logger.Errorf("invalid mattermost-url - expected format: 'http://<HOST>:<PORT>' - error: %s", err.Error())
		os.Exit(1)
	}

	mailServer := mail.New(*mailHost, *mailUser, *mailPass)

	fwdMappings, err := parseFwdMappings(*forward)
	if err != nil {
		logger.Errorf("unable to parse flag 'forward'. error: %s", err.Error())
		os.Exit(1)
	}

	if mailSubjectTemplate, err = template.New("mail-subject").Parse(*mailSubject); err != nil {
		logger.Errorf("invalid template for mail-subject - error: %s", err.Error())
		os.Exit(1)
	}
	if mailBodyTemplate, err = template.New("mail-body").Parse(*mailBody); err != nil {
		logger.Errorf("invalid template for mail-body - error: %s", err.Error())
		os.Exit(1)
	}

	//
	// "main loop"
	//
	//   * connect to mattermost
	//   * call 'dispatch' with forwards the messages if their contains a marker
	//   * if an error orccurs, reconnect to the mattermost server
	logger.Infof("startup - matterbot: v%s", version)
	for {
		logger.Info("connect to chat-server ...")
		chatServer, err := chat.Connect(url, *mattermostUser, *mattermostPass)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(2 * time.Second)
		} else {
			logger.Info("connected to chatServer")
			if err := dispatch(chatServer, mailServer, fwdMappings); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

// fwdMapping contains a pair of a marker and a corresponding mail-address.
type fwdMapping struct {
	marker   string
	mailAddr string
}

func parseFwdMappings(s string) ([]fwdMapping, error) {
	if len(s) == 0 {
		return nil, errors.New("flag 'forward' are mandatory")
	}

	fwdMappings := []fwdMapping{}
	for _, mapping := range strings.Split(s, ",") {
		x := strings.Split(mapping, "=")
		if len(x) != 2 {
			msg := "invalid format in flag 'forward': '%s' - valid example: 'user=abc@mail.com'"
			return nil, fmt.Errorf(msg, mapping)
		}
		marker := strings.TrimSpace(x[0])
		mailAddr := strings.TrimSpace(x[1])
		logger.Debugf("forward messages with marker: '@%s' to %s", marker, mailAddr)
		fwdMappings = append(fwdMappings, fwdMapping{marker, mailAddr})
	}
	return fwdMappings, nil
}
