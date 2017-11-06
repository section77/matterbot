// Package mail is the interface to the mail-system
package mail

import (
	"net/smtp"

	"github.com/section77/matterbot/logger"
)

// Server defines the interface to the mail-system
type Server interface {
	Send(*Message) error
}

// Header representes the mail-header
type Header struct {
	From      string
	To        string
	Subject   string
	Timestamp string
}

// Message represents an mail-message
type Message struct {
	Header  Header
	Body    string
	Content string
}

func New(host, user, pass string) Server {
	return &serverImpl{
		host: host,
		user: user,
		pass: pass,
	}
}

type serverImpl struct {
	host string
	user string
	pass string
}

// Send the given message
func (s *serverImpl) Send(msg *Message) error {
	auth := smtp.PlainAuth(
		"",
		s.user,
		s.pass,
		s.host,
	)

	logger.Info("send mail ...")
	return smtp.SendMail(
		s.host,
		auth,
		"matterbot@section77.de",
		[]string{"matterbot@j-keck.net"},
		[]byte(msg.Body),
	)
}
