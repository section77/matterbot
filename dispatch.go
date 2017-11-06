package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/section77/matterbot/chat"
	"github.com/section77/matterbot/logger"
	"github.com/section77/matterbot/mail"
)

// the dispatch function listens for new chat messages and forwards the messages
// to the mailing-list if their start with the '@ml' prefix
//
//   - dispatch block's until a error occurs
//   - if the message can't be fowarded to the mailing-list, the mail-server error
//     message are send as a reply to the original message in the chat-system
func dispatch(chatServer chat.Server, mailServer mail.Server) error {
	msgC, errC, err := chatServer.Listen()
	if err != nil {
		return err
	}

	logger.Info("waiting for chat messages to broadcast ...")
	for {
		select {
		case msg := <-msgC:
			content := strings.TrimSpace(msg.Content)

			// marker prefixes - every message which starts which any of this prefixes are broadcasted
			//  TODO: make this configurable? mapping of prefix -> receiver?
			//    - @ml-orga -> orga@....
			//    - @ml-intern -> mitglieder@...
			prefixes := []string{"@ml ", "@ml,"}

			if prefix, found := prefix(content, prefixes); found {
				logger.Infof("prefix: '%s' found, broadcast chat message to ml - from: %s in channel: %s",
					prefix, msg.UserName, msg.ChannelName)

				// remove the marker prefix
				msg.Content = strings.TrimLeft(content, prefix)

				// send the mail
				if err = mailServer.Send(composeMessageFromChat(&msg)); err != nil {
					logger.Errorf("unable to send mail - notify user in chat - mail error: %s", err.Error())
					if err = chatServer.Send(&chat.Message{
						ReplyToID:   msg.ID,
						ChannelID:   msg.ChannelID,
						ChannelName: msg.ChannelName,
						Content:     "matterbot error: " + err.Error(),
					}); err != nil {
						logger.Errorf("unable to notify user about mail error - i give up - sorry! - chat error: %s",
							err.Error())
					}
				}

			} else {
				logger.Debugf("ignore message from: '%s' - didn't contain the '@ml' prefix", msg.UserName)
			}
		case chatErr := <-errC:
			return chatErr
		}
	}
}

// if the given string starts with any of the given prefixes, return the found prefix
//
//   returns:
//     - found:     ("<PREFIX>", true)
//     - not found: ("", false)
func prefix(s string, ps []string) (string, bool) {
	for _, p := range ps {
		if strings.HasPrefix(s, p) {
			return p, true
		}
	}
	return "", false
}

// compose a mail message from the given chat message
func composeMessageFromChat(msg *chat.Message) *mail.Message {
	subject := fmt.Sprintf("Broadcast von chat.section77.de - '%s' schrieb in '%s'", msg.UserName, msg.ChannelName)
	// time format (https://tools.ietf.org/html/rfc5322#section-3.3)
	tsFmt := "Mon, 02 Jan 2006 15:04:05 MST"
	return mail.ComposeMessage(mail.Header{
		From:      "matterbot@section77.de",
		To:        "j@j-keck.net",
		Subject:   subject,
		Timestamp: time.Now().Format(tsFmt),
	}, msg.Content)
}
