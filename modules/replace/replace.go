package replace

import (
    "../../types"
    "regexp"
)

var lastMessage string

func Handle (message types.Message) types.Response {
    var response = types.Response{}

    regex := regexp.MustCompile(`^s\/(.*?)\/(.*?)(\/(.+))?$`)
    matches := regex.FindStringSubmatch(message.Message)

    if (len(matches) >= 2) {
        search := matches[1]
        replace := matches[2]
        // user := matches[2]

        replaceRegex := regexp.MustCompile(search)
        reply := replaceRegex.ReplaceAllString(lastMessage, replace)
        response.Messages = []string{reply}
    } else {
        lastMessage = message.Message
    }

    return response
}