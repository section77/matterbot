package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/section77/matterbot/chat"
	"github.com/section77/matterbot/logger"
	"github.com/section77/matterbot/mail"
)

// the dispatch function listens for new chat messages and forwards the messages
// per mail if their start with a special marker
//
//   - dispatch block's until a error occurs
//   - if the message can't be fowarded to per mail, the mail-server error
//     message are send as a reply to the original message in the chat-system
func dispatch(chatServer chat.Server, mailServer mail.Server, fwdMappings []fwdMapping) error {
	msgC, errC, err := chatServer.Listen()
	if err != nil {
		return err
	}

	logger.Info("waiting for chat messages to broadcast ...")
	for {
		select {
		case msg := <-msgC:

			if mappings, content, found := findMappings(msg.Content, fwdMappings); found {
				logger.Infof("%d marker found - chat-msg from: %s, in channel: %s - forward to each recipient",
					len(mappings), msg.UserName, msg.ChannelName)

				for _, m := range mappings {
					logger.Infof("forward message with marker: '%s' to %s", m.marker, m.mailAddr)

					// send the mail
					if err = mailServer.Send(composeMessage(&msg, content, m.mailAddr), *mailUseTLS); err != nil {
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
					} else {
						logger.Debugf("mail to %s delivered", m.mailAddr)
					}
				}
			} else {
				logger.Debugf("ignore message from: '%s' - didn't contain any configured marker", msg.UserName)
			}
		case chatErr := <-errC:
			return chatErr
		}
	}
}

// find all mappings from the given content
func findMappings(content string, allFwdMappings []fwdMapping) ([]fwdMapping, string, bool) {
	foundFwdMappings := []fwdMapping{}

	isSeparator := func(c rune) bool {
		return c == ' ' || c == ','
	}
	splitAt := func(s string, n int) (string, string) {
		if n > 0 && len(s) > n {
			return s[:n], s[n:]
		}
		return s, ""
	}

	var foundMarker string

	// we mutate this 'work' variable in each loop to remove and '@xxx' marker
	work := strings.TrimLeftFunc(content, unicode.IsSpace)
	for strings.HasPrefix(work, "@") {
		foundMarker, work = splitAt(work, strings.IndexFunc(work, isSeparator))

		// remove the separator from the content
		work = strings.TrimLeftFunc(work, isSeparator)

		// is the marker defined in 'allFwdMappings', add it to 'foundFwdMappings'
		for _, m := range allFwdMappings {
			if foundMarker == "@"+m.marker {
				foundFwdMappings = append(foundFwdMappings, m)
			}
		}
	}

	return foundFwdMappings, work, len(foundFwdMappings) > 0
}

// compose a mail message from a chat message

//   * meta-data are used from the given chat-message
//   * mail-content are used from the given 'content' paramter
func composeMessage(msg *chat.Message, content string, to string) *mail.Message {
	// TODO: mail subject configurable with placeholders like '$channel$, $user$, ...)
	subject := fmt.Sprintf("'%s' writes in mattermost channel: '%s'", msg.UserName, msg.ChannelName)

	// time format (https://tools.ietf.org/html/rfc5322#section-3.3)
	tsFmt := "Mon, 02 Jan 2006 15:04:05 MST"
	return mail.ComposeMessage(mail.Header{
		From:      *mailUser,
		To:        to,
		Subject:   subject,
		Timestamp: time.Now().Format(tsFmt),
	}, content)
}
