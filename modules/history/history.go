package history

import (
    "../../types"
    "fmt"
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

func Handle (message *types.Message) types.Response {
    var response types.Response

    if message.Command == "seen" {
        last, error := GetUserLatest(message.Arguments[0])
        if error == "" {
            response.Message = fmt.Sprintf("%s was seen in %s at %s: %s", last.Nick, last.Channel, last.Timestamp, last.Original)
        } else {
            response.Message = "{red}" + error
        }
    }


    // Don't record private messages or commands
    if message.Private || message.Replied {
        return response
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

    return response
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