package example

import (
    "../../types"
    "../../modules"
)

func init () {
    modules.RegisterModule("example", Handle)
}

func Handle (message *types.Message) types.Response {
    var response = types.Response{}

    response.Message = echo(message.Message)

    return response
}

func echo (message string) string {
    return message
}