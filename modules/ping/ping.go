package example

import (
    "../../types"
    "../../modules"
)

func init () {
    modules.RegisterModule("ping", Handle)
}

func Handle (message *types.Message) types.Response {
    var response = types.Response{}

    if message.Command == "ping" {
        response.Message = "Pong"
    }

    return response
}