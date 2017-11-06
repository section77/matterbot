package mail

import (
	"bytes"
)

// ComposeMessage composes an mail-message from the given
// mail-header and content
func ComposeMessage(header Header, content string) *Message {
	mcb := newMessageContentBuilder()
	mcb.AppendHeader("From", header.From)
	mcb.AppendHeader("To", header.To)
	mcb.AppendHeader("Subject", header.Subject)
	mcb.AppendHeader("Date", header.Timestamp)
	mcb.AppendContent(content)

	return &Message{header, mcb.String(), content}
}

type messageContentBuilder struct {
	buf bytes.Buffer
}

func newMessageContentBuilder() messageContentBuilder {
	return messageContentBuilder{}
}

func (mb *messageContentBuilder) AppendHeader(n, v string) {
	mb.buf.WriteString(n)
	mb.buf.WriteString(": ")
	mb.buf.WriteString(v)
	mb.buf.WriteString("\r\n")
}

func (mb *messageContentBuilder) AppendContent(s string) {
	// header / content seperator
	mb.buf.WriteString("\r\n")
	mb.buf.WriteString(s)
}

func (mb *messageContentBuilder) String() string {
	return mb.buf.String()
}
