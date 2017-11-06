package chat

import (
	"time"

	"github.com/section77/matterbot/logger"
)

// ServerMock implements the chat-system interface to use it in unit-tests
type ServerMock struct {
	connected bool

	Messages []*Message

	msgC chan Message
	errC chan error
}

// NewMock instantiates a new ServerMock
func NewMock() *ServerMock {
	return &ServerMock{
		connected: true,
		msgC:      make(chan Message, 100),
		errC:      make(chan error, 1),
	}
}

// IsConnected returns the current connection status which is controlled in 'ServerMock.connected'
func (mock *ServerMock) IsConnected() bool {
	return mock.connected
}

// Send emulates an send-action and saves all messages in the mock// Send emulates an send-action and saves all messages in the mock
func (mock *ServerMock) Send(msg *Message) error {
	logger.Debugf("send per chat in channel: %s message: %s", msg.ChannelName, msg.Content)
	mock.Messages = append(mock.Messages, msg)
	return nil
}

// Listen returns a channel with chat messages and one with error messages.
//  * chat messages can be triggered per 'TriggerMsgEvent'
//  * error events can be triggered per 'TriggerErrorevent'
// on this mock
func (mock *ServerMock) Listen() (<-chan Message, <-chan error, error) {
	return mock.msgC, mock.errC, nil
}

// TriggerMsgEvent triggers an event in the 'Message channel'
// which is returned from the 'Listen' function
func (mock *ServerMock) TriggerMsgEvent(msg Message) {
	mock.msgC <- msg

	// give the 'dispatcher' time to react
	time.Sleep(100 * time.Millisecond)
}

// TriggerErrorEvent triggers an event in the 'error channel'
// which is returned from the 'Listen' function
func (mock *ServerMock) TriggerErrorEvent(err error) {
	mock.errC <- err

	// give the 'dispatcher' time to react
	time.Sleep(100 * time.Millisecond)
}
