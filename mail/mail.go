// Package mail is the interface to the mail-system
package mail

import (
	"crypto/tls"
	"io"
	"net"
	"net/smtp"

	"github.com/section77/matterbot/logger"
)

// Server defines the interface to the mail-system
type Server interface {
	Send(*Message, bool) error
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
func (s *serverImpl) Send(msg *Message, useTLS bool) error {
	logger.Debugf("send mail (per %s) - host: %s, from: %s, to: %s",
		protocolStr(useTLS), s.host, msg.Header.From, msg.Header.To)

	host, _, _ := net.SplitHostPort(s.host)
	auth := smtp.PlainAuth(
		"",
		s.user,
		s.pass,
		// the host part can't have a port number
		host,
	)

	if useTLS {
		return s.sendPerTLS(host, auth, msg)
	}
	return s.sendPerSTARTTLS(auth, msg)
}

func protocolStr(useTLS bool) string {
	protocol := "STARTTLS"
	if useTLS {
		protocol = "TLS"
	}
	return protocol
}

func (s *serverImpl) sendPerSTARTTLS(auth smtp.Auth, msg *Message) error {
	return smtp.SendMail(
		s.host,
		auth,
		msg.Header.From,
		[]string{msg.Header.To},
		[]byte(msg.Body),
	)
}

func (s *serverImpl) sendPerTLS(host string, auth smtp.Auth, msg *Message) error {
	var err error
	var con *tls.Conn
	var client *smtp.Client
	var writer io.WriteCloser

	if con, err = tls.Dial("tcp", s.host, &tls.Config{
		ServerName: host,
	}); err != nil {
		return err
	}

	if client, err = smtp.NewClient(con, host); err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(msg.Header.From); err != nil {
		return err
	}

	if err = client.Rcpt(msg.Header.To); err != nil {
		return err
	}

	if writer, err = client.Data(); err != nil {
		return err
	}

	if _, err = writer.Write([]byte(msg.Body)); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}
