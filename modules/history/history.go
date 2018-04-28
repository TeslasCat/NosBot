package history

import (
    "../../types"
)


type Channel struct {
    History []*types.Message
}

type User struct {
    Message *types.Message
}

var channels map[string]*Channel
var users map[string]*User


func init() {
    channels = make(map[string]*Channel)
    users = make(map[string]*User)
}

func Handle (message *types.Message) {
    // Don't record private messages or commands
    if message.Private || message.Replied {
        return
    }

    // Update channel history
    if _, exists := channels[message.Channel]; !exists {
        channels[message.Channel] = &Channel{}
    }
    channels[message.Channel].History = append(channels[message.Channel].History, message)

    // Update user history
    if _, exists := users[message.Nick]; !exists {
        users[message.Nick] = &User{}
    }
    users[message.Nick].Message = message
}

func GetUserLatest(nick string) (*types.Message, string) {
    if _, exists := users[nick]; !exists {
        return nil, "Nick not found"
    }

    return users[nick].Message, ""
}

func GetChannelLatest(channel string) (*types.Message, string) {
    if _, exists := channels[channel]; !exists {
        return nil, "No channel history found"
    }

    history := channels[channel].History

    return history[len(history)-1], ""
}