// Package chat is the interface to the chat-system
package chat

// Server defines the interface to the chat-system
type Server interface {
	IsConnected() bool
	Send(*Message) error
	Listen() (<-chan Message, <-chan error, error)
}

// Message represents a chat message
type Message struct {
	ID          string
	UserID      string
	UserName    string
	ChannelID   string
	ChannelName string
	Content     string
	ReplyToID   string
}
