package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/section77/matterbot/chat"
	"github.com/section77/matterbot/mail"
)

// only messages with a special prefix should be forwarded per mail
func TestDispatchFiltersMessagesPerPrefix(t *testing.T) {
	chatMock := chat.NewMock()
	mailMock := mail.NewMock()

	go dispatch(chatMock, mailMock)

	tests := []struct {
		content              string
		msgShouldBeForwarded bool
	}{
		{"without prefix", false},
		{"@ml with prefix", true},
		{"@mlwith prefix but without a space", false},
		{"@ml,also with prefix", true},
	}

	// send all test messages
	for _, test := range tests {
		chatMock.TriggerMsgEvent(chat.Message{
			UserName:    "test-filter",
			ChannelName: "test-filter",
			Content:     test.content,
		})
	}

	if len(mailMock.Messages) != 2 {
		t.Errorf("expected 2 mail messages, but found: %d messages", len(mailMock.Messages))
	}

	var n int
	for _, test := range tests {
		if test.msgShouldBeForwarded {
			content := mailMock.Messages[n].Content
			if strings.HasSuffix(content, test.content) {
				t.Errorf("mail content from the %d. mail should end with '%s', but ends with '%s'",
					n+1, test.content, content)
			}
			n = n + 1
		}
	}

}

// the call on 'dispatch' should block, and only returns
// if a error occurs
func TestDispatchBlocksAndReturnsTheError(t *testing.T) {
	chatMock := chat.NewMock()
	mailMock := mail.NewMock()

	go func() {
		time.Sleep(800 * time.Millisecond)
		chatMock.TriggerErrorEvent(errors.New("test-error"))
	}()

	blockingStartTs := time.Now()
	err := dispatch(chatMock, mailMock)

	if time.Since(blockingStartTs) < 400*time.Millisecond {
		t.Errorf("'dispatch' call didn't block")
	}

	if err == nil || err.Error() != "test-error" {
		t.Errorf("expected error not received - received: '%+v'", err)
	}
}

// if the bot can't forward a chat-message per mail, the error
// should be forwarded in the chat
func TestDispatcherSendsChatMsgOnMailError(t *testing.T) {
	chatMock := chat.NewMock()
	mailMock := mail.NewMock()

	go dispatch(chatMock, mailMock)

	dummyMsg := chat.Message{
		UserName:    "test-user",
		ChannelName: "test-channel",
		Content:     "@ml dummy message",
	}

	//
	// validate the 'good path'
	//   - send a chat message with @ml prefix
	//   - the bot should not send any chat messages
	//
	chatMock.TriggerMsgEvent(dummyMsg)
	if len(chatMock.Messages) != 0 {
		t.Error("expected empty chat-message queue")
	}

	//
	// validate the 'bad path'
	//   - configure mail-mock to return a error
	//   - send a chat message with @ml prefix
	//   - dispatcher should send a chat message with the mail-server error
	//
	mailMock.SetMailServerError(errors.New("mail-mock-test-error"))
	chatMock.TriggerMsgEvent(dummyMsg)
	if len(chatMock.Messages) != 1 {
		t.Errorf("expected one message in chat-mock, actual message count in chat-mock: %d", len(chatMock.Messages))
	}

	if !strings.Contains(chatMock.Messages[0].Content, "mail-mock-test-error") {
		t.Errorf("expected message in chat-mock not found")
	}
}
