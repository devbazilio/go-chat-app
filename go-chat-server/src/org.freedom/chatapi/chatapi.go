package chatapi

import (
	"net/http"
	"org.freedom/bootstrap"
	"sync"
)

var wsHandlers = bootstrap.HttpHandler{
	ApiHandlers: map[string]bootstrap.ApiHandler{
		"get": wsHandler,
	},
}

type ChannelsList struct {
	mutex    sync.Mutex
	channels map[string]bool
}

type channelsJSON struct {
	Channels *[]string `json:"channels"`
}

type messagesJSON struct {
	Messages *[]channelMessage `json:"messages"`
}

var channelsList = ChannelsList{
	channels: map[string]bool{
		"general": true,
		"news":    true,
	},
}

type channelMessage struct {
	Time    int64  `json:"time"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

var channelMessages = make(map[string][]channelMessage)

func Setup() {
	bootstrap.AddEndPoints("/ws", &wsHandlers)
	bootstrap.AddCommandListener("SET_USERNAME", commandSetUserName)
	bootstrap.AddCommandListener("GET_CHANNELS", commandListChannels)
	bootstrap.AddCommandListener("GET_CHANNEL_MESSAGES", commandListChannelMessages)
	bootstrap.AddCommandListener("POST_MESSAGE", commandStoreUserMessage)
}

/*
func addChannel(r *http.Request) (status int, response *[]byte, e error) {
	name := r.FormValue("name")
	if len(name) > 0 && len(name) < 255 {
		channelsList.mutex.Lock()
		channelsList.channels[name] = true
		channelsList.mutex.Unlock()
		return listChannels(r)
	}
	return http.StatusBadRequest, nil, nil
}


*/
func wsHandler(r *http.Request) (status int, response *[]byte, e error) {
	var body = []byte("PONG")
	return http.StatusOK, &body, nil
}
