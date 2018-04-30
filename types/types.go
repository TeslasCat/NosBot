package types

// https://mholt.github.io/json-to-go/ <3

type Config struct {
    Server         string   `json:"server"`
    Channels       []string `json:"channels"`
    Nick           string   `json:"nick"`
    User           string   `json:"user"`
    // Nickserv       string   `json:"nickserv"`
    Debug          bool     `json:"debug"`
    Port           int      `json:"port"`
    Secure         bool     `json:"secure"`
    SkipVerify     bool     `json:"skipVerify"`
    Admin          []string `json:"admin"`
    // WordnikAPI     string   `json:"wordnik_api"`
    // Greeting       []string `json:"greeting"`
    // GreetingIgnore []string `json:"greeting-ignore"`
    Modules       []string `json:"modules"`
    MatrixUser      string `json:"matrixUser"`
    MatrixPassword  string `json:"matrixPassword"`
}

type Message struct {
    Nick string
    Channel string
    Message string
    Original string
    Timestamp string
    Private bool
    Command string
    Arguments []string
    Replied bool
    Platform string
}

type Response struct {
    Type string     // action, message
    Target string
    Message string
    Messages []string
}