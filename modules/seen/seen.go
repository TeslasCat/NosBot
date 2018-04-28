package seen

import (
    "../../types"
    "../../modules"
    "../history"
    "fmt"
)

func init () {
    modules.RegisterModule("seen", Handle)
}

func Handle (message *types.Message) types.Response {
    var response types.Response

    if message.Command != "seen" {
        return response
    }

    if len(message.Arguments) > 0 {
        nick := message.Arguments[0]

        last, error := history.GetUserLatest(nick)
        if error == "" {
            response.Message = fmt.Sprintf("%s was seen in %s at %s: %s", last.Nick, last.Channel, last.Timestamp, last.Original)
        } else {
            response.Type = "action"
            response.Message = fmt.Sprintf("{red} cant see %s anywhere", nick)
        }
    } else {
        response.Message = "You will need to be more specific"
    }

    return response
}