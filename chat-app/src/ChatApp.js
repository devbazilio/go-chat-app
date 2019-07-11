import React from 'react';
import PropTypes from 'prop-types';
import {serverAddress} from './App';
import ChatDialogue from './ChatDialogue';

import {Sidebar} from './items';

export const DataContext = React.createContext({});

class ChatApp extends React.Component {
    state = {
        activeChannel: null,
        isPrivate: false,
        connected: false,
        channels: [],
        users: {},
        unreadChannels: {},
    };

    static propTypes = {
        userName: PropTypes.string.isRequired
    };

    constructor(props) {
        super(props);
        this.dialogueCallback = null;
    }

    componentDidMount() {
        this.openServerConnection();
    }

    openServerConnection = () => {
        const socket = new WebSocket(serverAddress);
        socket.onopen = this.onConnectionOpen;

        socket.onclose = event => {
            let reconnect = false;
            if (event.wasClean) {
                console.log('WS:DISCONNECTED');
            } else {
                reconnect = true;
                console.log('WS:DISCONNECTED (abort)');
            }
            console.log(`WS:DISCONNECTED (${event.code} ${event.reason})`);
            this.onConnectionClose(reconnect);
        };

        socket.onmessage = event => {
            console.log('DATA', event.data);
            try {
                const data = JSON.parse(event.data);
                if (data) {
                    this.onServerData(data);
                }
            } catch (e) {

            }
        };

        socket.onerror = function (error) {
            console.log('WS:ERROR', error);
        };
        this.socket = socket;
    };

    onConnectionOpen = () => {
        console.log('WS:OK');
        this.setState({connected: true});
        this.timeoutID = null;


        this.setName();
    };

    onConnectionClose = reconnect => {
        this.setState({connected: false});
        this.socket = null;
        if (reconnect) {
            this.timeoutID = setTimeout(this.openServerConnection, 10000);
        }
    };

    getChannels = () => this.sendCommand('GET_CHANNELS', null);
    setName = () => this.sendCommand('SET_USERNAME', this.props.userName);
    createChannel = (channel, isPrivate) => this.sendCommand('CREATE_CHANNEL', {channel, isPrivate});

    sendCommand = (command, data) => this.socket && this.socket.send(JSON.stringify({data, command}));

    onServerData = data => {
        const newState = {};
        const {activeChannel, unreadChannels} = this.state;

        if (!this.socket) {
            return;
        }

        if (data.channels) {
            newState.channels = data.channels;
            if (!activeChannel) {
                newState.activeChannel = {
                    name: data.channels[0],
                    isPrivate: false
                };
            }
        }

        if (this.dialogueCallback) {
            (data.messages || data.message) && this.dialogueCallback(data);
        }

        if (data.message && activeChannel && data.channelName !== activeChannel.name) {
            newState.unreadChannels = {
                ...unreadChannels,
                ...{[data.channelName]: true}
            }
        }

        if (data.users) {
            newState.users = data.users;
        }

        if (Object.keys(newState).length) {
            this.setState(newState);
        }
    };

    componentWillUnmount() {
        this.socket.close();
        this.socket = null;
        if (this.timeoutID) {
            clearTimeout(this.timeoutID);
        }
    }

    setActiveChannel = (name, isPrivate) => {
        const unreadChannels = {...this.state.unreadChannels};
        delete unreadChannels[name];

        this.setState({
            activeChannel: {name, isPrivate},
            isPrivate,
            unreadChannels
        });
    };

    setDialogueCallback = callback => {
        this.dialogueCallback = callback;
    };

    loadMessages = () => this.sendCommand('GET_CHANNEL_MESSAGES', this.state.activeChannel.name);
    getUsersList = () => this.sendCommand('LIST_USERS', null);

    sendUserMessage = message => this.sendCommand('POST_MESSAGE', {channel: this.state.activeChannel.name, message});

    askForChannelName = e => {
        e.preventDefault();
        e.stopPropagation();
        const channel = window.prompt('Type a channel name');
        if (channel && channel.trim().length) {
            this.createChannel(channel, false);
        }
    };

    render() {
        const {userName} = this.props;
        const {channels, connected, activeChannel, isPrivate, unreadChannels, users} = this.state;
        const contextData = {
            userName, connected, channels, activeChannel, unreadChannels, users,
            askForChannelName: this.askForChannelName,
            setActiveChannel: this.setActiveChannel,
            getUsersList: this.getUsersList,
        };

        return (
            <DataContext.Provider
                value={contextData}>
                <Sidebar/>
                {
                    activeChannel &&
                    <ChatDialogue key={Object.keys(activeChannel).join()}
                                  activeChannel={activeChannel}
                                  isPrivate={isPrivate}
                                  setCallback={this.setDialogueCallback}
                                  sendUserMessage={this.sendUserMessage}
                                  loadMessages={this.loadMessages}
                    />
                }
            </DataContext.Provider>
        )
    }
}

export default ChatApp