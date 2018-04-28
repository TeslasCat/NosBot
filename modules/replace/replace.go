package replace

import (
    "../../types"
    "../../modules"
    "regexp"
    "../history"
)

func init () {
    modules.RegisterModule("replace", Handle)
}

func Handle (message *types.Message) types.Response {
    var response = types.Response{}

    regex := regexp.MustCompile(`^s\/(.*?)\/(.*?)(?:\/(.+))?$`)
    matches := regex.FindStringSubmatch(message.Message)

    if (len(matches) >= 2) {
        search := matches[1]
        replace := matches[2]
        user := matches[3]

        var subject *types.Message
        var error string

        if user == "" {
            subject, error = history.GetChannelLatest(message.Channel)
        } else {
            subject, error = history.GetUserLatest(user)
        }

        if error != "" {
            return response
        }

        replaceRegex := regexp.MustCompile(search)
        reply := replaceRegex.ReplaceAllString(subject.Original, replace)
        response.Messages = []string{reply}
    }

    return response
}