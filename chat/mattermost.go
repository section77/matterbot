package chat

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/mattermost/platform/model"
	"github.com/section77/matterbot/logger"
)

// Mattermost implements the chat-system interface to use it with the 'mattermost' system
type Mattermost struct {
	client *model.Client4
	userID string
}

// Connect to the mattermost server.
// returns a connection handle to interact with the server
func Connect(url *url.URL, loginID, pass string) (*Mattermost, error) {
	client := model.NewAPIv4Client(url.String())

	logger.Debugf("try to login with loginID: '%s'", loginID)
	user, resp := client.Login(loginID, pass)
	if resp.Error != nil {
		err := fmt.Errorf("login error for loginID: '%s': %s", loginID, detailedErrOrMsg(resp))
		return nil, err
	}
	logger.Debugf("login success for loginID: %s, user: %+v", loginID, user)

	return &Mattermost{
		client: client,
		userID: user.Id,
	}, nil
}

// IsConnected returns the connection status.
func (m *Mattermost) IsConnected() bool {
	if msg, _ := m.client.GetPing(); msg == "OK" {
		return true
	}
	return false
}

// Send sends the given message.
// the 'ReplyToID' field from the message are used to respond to an message
func (m *Mattermost) Send(msg *Message) error {
	logger.Debugf("send msg in channel: %s, msg: %s", msg.ChannelName, msg.Content)
	post := &model.Post{
		RootId:    msg.ReplyToID,
		ChannelId: msg.ChannelID,
		Message:   msg.Content,
	}

	if _, resp := m.client.CreatePost(post); resp.Error != nil {
		err := fmt.Errorf("unable to post message in channel: %s, msg: %s, error: %s",
			msg.ChannelName, msg.Content, detailedErrOrMsg(resp))
		return err
	}

	return nil
}

// Listen connects to the mattermost server and listens for new messages.
//
// If the connection to the mattermost server fails, the error are returned.
//
// The two returned channels represents:
//
//   * an message channel with incomming chat messages
//   * an error channel with occurred errors
//
func (m *Mattermost) Listen() (<-chan Message, <-chan error, error) {

	// url is already validated - so no error checking here
	url, _ := url.Parse(m.client.Url)
	wsURL := fmt.Sprintf("ws://%s:%s", url.Hostname(), url.Port())

	wsClient, wsErr := model.NewWebSocketClient(wsURL, m.client.AuthToken)
	if wsErr != nil {
		err := fmt.Errorf("ws connection error: %s", wsErr.Error())
		return nil, nil, err
	}

	msgC := make(chan Message, 100)
	errC := make(chan error)
	go func() {
		wsClient.Listen()
		defer wsClient.Close()

		for {
			event := <-wsClient.EventChannel
			if event == nil {
				errC <- errors.New("'nil' event received - disconnected?")
				return
			} else if event.Event == model.WEBSOCKET_EVENT_POSTED {
				if post := model.PostFromJson(strings.NewReader(event.Data["post"].(string))); post != nil {
					userName := "id:" + post.UserId
					if user, err := m.GetUser(post.UserId); err == nil {
						userName = user.Username
					}

					channelName := "id:" + post.ChannelId
					if channel, err := m.GetChannel(post.ChannelId); err == nil {
						channelName = channel.Name
					}

					msg := Message{
						ID:          post.Id,
						UserID:      post.UserId,
						UserName:    userName,
						ChannelID:   post.ChannelId,
						ChannelName: channelName,
						Content:     post.Message,
					}
					logger.Debugf("publish new message from: '%s', in channel: '%s'", userName, channelName)
					msgC <- msg
				}
			}
		}
	}()

	return msgC, errC, nil
}

func (m *Mattermost) GetUser(userID string) (*model.User, error) {
	logger.Debugf("try to lookup user by id: '%s'", userID)

	etag := ""
	user, resp := m.client.GetUser(userID, etag)
	if resp.Error != nil {
		err := fmt.Errorf("user with id: '%s' not found: %s", userID, detailedErrOrMsg(resp))
		return nil, err
	}

	logger.Debugf("user with id: '%s' found, user: %+v", userID, user)
	return user, nil
}

func (m *Mattermost) GetChannel(channelID string) (*model.Channel, error) {
	logger.Debugf("try to lookup channel by id: '%s'", channelID)

	etag := ""
	channel, resp := m.client.GetChannel(channelID, etag)
	if resp.Error != nil {
		err := fmt.Errorf("channel with id: '%s' not found: %s", channelID, detailedErrOrMsg(resp))
		return nil, err
	}

	logger.Debugf("channel with id: '%s' found, channel: %+v", channelID, channel)
	return channel, nil
}

func (m *Mattermost) GetTeamByName(name string) (*model.Team, error) {
	logger.Debugf("try to lookup team by name: '%s'", name)

	etag := ""
	team, resp := m.client.GetTeamByName(name, etag)
	if resp.Error != nil {
		err := fmt.Errorf("team with name: '%s' not found: %s", name, detailedErrOrMsg(resp))
		return nil, err
	}

	logger.Debugf("team with name: '%s' found, team: %+v", name, team)
	return team, nil
}

func (m *Mattermost) GetChannelByName(name string, team *model.Team) (*model.Channel, error) {
	logger.Debugf("try to lookup channel by name: %s, in team: %s", name, team.Name)

	etag := ""
	channel, resp := m.client.GetChannelByName(name, team.Id, etag)
	if resp.Error != nil {
		err := fmt.Errorf("channel with name: '%s' in team: '%s' not found: %s", name, team.Name, detailedErrOrMsg(resp))
		return nil, err
	}

	logger.Debugf("channel with name: '%s' in team: '%s' found, channel: %+v", name, team.Name, channel)
	return channel, nil
}

// try to get the detailed error message from the response.
// if it's empty, return the general error message
func detailedErrOrMsg(resp *model.Response) string {
	if resp.Error.DetailedError != "" {
		return resp.Error.DetailedError
	}
	return resp.Error.Message
}
