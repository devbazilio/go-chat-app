package chatapi

import (
	"net/http"
	"org.freedom/bootstrap"
	"sync"
	"time"
)

var wsHandlers = bootstrap.HttpHandler{
	ApiHandlers: map[string]bootstrap.ApiHandler{
		"get": wsHandler,
	},
}

type ChannelsList struct {
	mutex    sync.Mutex
	channels map[string]*channelPeers
}

type channelJSON struct {
	IsCommon bool `json:"is_common"`
}

type channelsJSON struct {
	Channels map[string]channelJSON `json:"channels"`
}

type messagesJSON struct {
	Messages *[]channelMessage `json:"messages"`
}

type messageJSON struct {
	ChannelName string         `json:"channelName"`
	Message     channelMessage `json:"message"`
}

type channelPeers struct {
	mutex    sync.Mutex
	isCommon bool
	peers    []string
}

var channelsList = ChannelsList{
	channels: make(map[string]*channelPeers),
}

type channelMessage struct {
	Time    int64  `json:"time"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

type channelsMessagesMap map[string][]channelMessage
type channelMessagesHistory struct {
	mutex    sync.Mutex
	messages channelsMessagesMap
}

var channelMessages = channelMessagesHistory{
	messages: make(channelsMessagesMap, 0),
}

func (c *channelMessagesHistory) AppendMessage(channelName, text, sender string) *channelMessage {
	var newMessage = channelMessage{
		Message: text,
		Time:    time.Now().Unix(),
		Sender:  sender,
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	channelsMessagesArray, _ := c.messages[channelName]
	if channelsMessagesArray == nil {
		channelsMessagesArray = make([]channelMessage, 0, 1)
	}
	channelsMessagesArray = append(channelsMessagesArray, newMessage)
	channelMessages.messages[channelName] = channelsMessagesArray
	return &newMessage
}

func (cl *ChannelsList) AddChannel(name string, isCommon bool) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	_, exists := cl.channels[name]
	if !exists {
		cl.channels[name] = &channelPeers{isCommon: isCommon}
	}
}

func Setup() {
	bootstrap.AddEndPoints("/ws", &wsHandlers)
	bootstrap.AddCommandListener("SET_USERNAME", commandSetUserName)
	bootstrap.AddCommandListener("GET_CHANNELS", commandListChannels)
	bootstrap.AddCommandListener("GET_CHANNEL_MESSAGES", commandListChannelMessages)
	bootstrap.AddCommandListener("POST_MESSAGE", commandStoreUserMessage)
	bootstrap.AddCommandListener("CREATE_CHANNEL", commandCreateChannel)
	channelsList.AddChannel("general", true)
	channelsList.AddChannel("news", true)
}

func wsHandler(r *http.Request) (status int, response *[]byte, e error) {
	var body = []byte("PONG")
	return http.StatusOK, &body, nil
}
