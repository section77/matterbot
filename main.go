package main

import (
	"net/url"
	"os"
	"time"

	"github.com/namsral/flag"

	"github.com/section77/matterbot/chat"
	"github.com/section77/matterbot/logger"
	"github.com/section77/matterbot/mail"
)

var (
	logVerbose  = flag.Bool("v", false, "enable verbose / debug output")
	logDisabled = flag.Bool("q", false, "disable logging / be quite")

	mattermostURL  = flag.String("mattermost-url", "http://127.0.0.1:8065", "mattermost url (http://x.x.x.x:xxxx)")
	mattermostUser = flag.String("mattermost-user", "matterbot", "mattermost user")
	mattermostPass = flag.String("mattermost-pass", "tobrettam", "mattermost password")

	mailHost = flag.String("mail-host", "127.0.0.1:25", "mail-server host (<HOST>:<PORT>)")
	mailUser = flag.String("mail-user", "matterbot", "mail login user")
	mailPass = flag.String("mail-pass", "tobrettam", "mail login pass")
)

func main() {
	logger.Infof("startup")
	flag.Parse()

	if *logVerbose {
		logger.SetLogLevel(logger.DebugLevel)
	} else if *logDisabled {
		logger.SetLogLevel(logger.Disabled)
	}

	url, err := url.Parse(*mattermostURL)
	if err != nil {
		logger.Errorf("invalid mattermost-url - expected format: 'http://<HOST>:<PORT>' - error: %s", err.Error())
		os.Exit(1)
	}

	mailServer := mail.New(*mailHost, *mailUser, *mailPass)

	for {
		logger.Info("connect to chat-server ...")
		chatServer, err := chat.Connect(url, *mattermostUser, *mattermostPass)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(2 * time.Second)
		} else {
			logger.Info("connected to chatServer")
			if err := dispatch(chatServer, mailServer); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
