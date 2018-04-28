package replace

import (
    "../../types"
    "regexp"
    "../history"
    "log"
)

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
            log.Print(user)
            subject, error = history.GetUserLatest(user)
        }

        if error != "" {
            return response
        }

        log.Printf("Replacing %s", subject.Original)

        replaceRegex := regexp.MustCompile(search)
        reply := replaceRegex.ReplaceAllString(subject.Original, replace)
        response.Messages = []string{reply}
    }

    return response
}