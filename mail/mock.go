package mail

import (
	"github.com/section77/matterbot/logger"
)

// ServerMock implements the mail-system interface to use it in unit-tests
type ServerMock struct {
	MailServerError error
	Messages        []*Message
}

// NewMock instantiates a new ServerMock
func NewMock() *ServerMock {
	return &ServerMock{}
}

// SetMailServerError set's the error which should be returned
// when the 'Send' function are called.
func (mock *ServerMock) SetMailServerError(err error) {
	mock.MailServerError = err
}

// Send emulates an send-action and saves all messages in the mock.
// If the 'SetMailServerError' are called with an error, this function
// returns the stored error.
func (mock *ServerMock) Send(msg *Message) error {
	if mock.MailServerError == nil {
		logger.Debugf("send per mail: %s", msg.Content)
		mock.Messages = append(mock.Messages, msg)
		return nil
	}

	logger.Debugf("mail-mock is configured to trigger an error - returning the error")
	return mock.MailServerError
}
