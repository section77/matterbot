package main

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/section77/matterbot/chat"
	"github.com/section77/matterbot/mail"
)

func TestFindFwdMappings(t *testing.T) {
	expectedMappings := []fwdMapping{
		fwdMapping{"user1", "user1@mail.com"},
		fwdMapping{"user2", "user2@mail.com"},
	}
	expectedContent := "test message"
	mappings, content, found := findMappings("@user1, @xx @user2 "+expectedContent, expectedMappings)

	// found should be true
	if !found {
		t.Errorf("found was 'false' - shout be 'true'")
	}

	// mappings should contain all expected mappings
	for i, m := range mappings {
		if m != expectedMappings[i] {
			t.Errorf("mapping don't match - expected: %+v, received: %+v", expectedMappings[i], m)
		}
	}

	// content should be 'test message'
	if content != expectedContent {
		t.Errorf("unexpected content: %s, exepected: %s", content, expectedContent)
	}
}

// only messages with a special marker should be forwarded per mail
func TestDispatchSendMessagesOnlyWithMarker(t *testing.T) {
	chatMock := chat.NewMock()
	mailMock := mail.NewMock()

	go dispatch(chatMock, mailMock, []fwdMapping{
		fwdMapping{"ml", "ml@mail.com"},
	})

	tests := []struct {
		content              string
		msgShouldBeForwarded bool
	}{
		{"without prefix", false},
		{"@ml with prefix", true},
		{"  @ml  with prefix and spaces", true},
		{"@mlwith prefix but without a space", false},
		{"@ml,also with prefix", true},
	}
	// grr: i need something like: 'len(tests.filter(_.msgShouldBeForwarded))
	expectedMessagesCount := 3

	// msgs all test messages
	for _, test := range tests {
		chatMock.TriggerMsgEvent(chat.Message{
			UserName:    "test",
			ChannelName: "test",
			Content:     test.content,
		})
	}

	if len(mailMock.Messages) != expectedMessagesCount {
		t.Errorf("expected %d mail messages, but found: %d messages",
			expectedMessagesCount, len(mailMock.Messages))
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

func TestDispatcherSendsMailToAllRecipients(t *testing.T) {
	chatMock := chat.NewMock()
	mailMock := mail.NewMock()

	go dispatch(chatMock, mailMock, []fwdMapping{
		fwdMapping{"user1", "user1@mail.com"},
		fwdMapping{"user2", "user2@mail.com"},
	})

	// one receiver
	chatMock.TriggerMsgEvent(chat.Message{
		Content: " @user1 hey",
	})
	verifyDispatchSendsMailToAllRecipients("one receiver", mailMock.Messages, []string{"user1@mail.com"}, t)
	mailMock.ClearMessages()

	// two receiver
	chatMock.TriggerMsgEvent(chat.Message{
		Content: "@user1,@user2 hey",
	})
	verifyDispatchSendsMailToAllRecipients("two receivers", mailMock.Messages, []string{"user1@mail.com", "user2@mail.com"}, t)
	mailMock.ClearMessages()
}

func verifyDispatchSendsMailToAllRecipients(name string, msgs []*mail.Message, expectedRecipients []string, t *testing.T) {
	if len(msgs) != len(expectedRecipients) {
		t.Errorf("%s: expected %d mail-messages, but %d messages received", name, len(expectedRecipients), len(msgs))
		return
	}

	for i, msg := range msgs {
		if !strings.HasPrefix(msg.Header.To, expectedRecipients[i]) {
			t.Errorf("%s: expected recipent: %s, found: %s", name, expectedRecipients[i], msg.Header.To)
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
	err := dispatch(chatMock, mailMock, []fwdMapping{})

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

	go dispatch(chatMock, mailMock, []fwdMapping{
		fwdMapping{"ml", "ml@mail.com"},
	})

	dummyMsg := chat.Message{
		Content: "@ml dummy message",
	}

	//
	// validate the 'good path'
	//   - msgs a chat message with @ml prefix
	//   - the bot should not msgs any chat messages
	//
	chatMock.TriggerMsgEvent(dummyMsg)
	if len(chatMock.Messages) != 0 {
		t.Error("expected empty chat-message queue")
	}

	//
	// validate the 'bad path'
	//   - configure mail-mock to return a error
	//   - msgs a chat message with @ml prefix
	//   - dispatcher should msgs a chat message with the mail-server error
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
