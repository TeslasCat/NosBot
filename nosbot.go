package main

import (
	"os"
	"log"
	"regexp"
	"strings"
	"encoding/json"
	"crypto/tls"
	"github.com/lrstanley/girc"
	"../mautrix-go"
	"flag"
	"fmt"
	"time"

	"./types"
	"./modules"
	"./modules/history"
	_ "./modules/notes"
	_ "./modules/replace"
	_ "./modules/seen"
)

var ircClient *girc.Client
var matrixClient *mautrix.MatrixBot

func main() {

	conf := loadConfig()

	if conf.Debug {
		log.Printf("Printing Configuration file: \n%+v\n", conf)
	}

	// Setup matrix client
	var homeserver = flag.String("homeserver", "https://tak.lward.co.uk", "Macak homeserver")
	var username = flag.String("username", conf.MatrixUser, "Matrix username localpart")
	var password = flag.String("password", conf.MatrixPassword, "Matrix password")

	matrixClient = mautrix.Create(*homeserver)
	err := matrixClient.PasswordLogin(*username, *password)
	if err != nil {
		panic(err)
	}

	// err = matrixClient.Join(conf.MatrixRoom)
	// if err != nil {
	// 	panic(err)
	// }

	stop := make(chan bool, 1)
	go matrixClient.Listen()
	go func() {
	Loop:
		for {
			select {
			case <-stop:
				break Loop
			case evt := <-matrixClient.Timeline:
				evt.MarkRead()
				switch evt.Type {
				case mautrix.EvtRoomMessage:
					if (evt.Type == "m.room.message") {
						log.Print(evt);

						message := types.Message{}
						message.Nick = evt.Sender
						message.Message = evt.Content["body"].(string)
						message.Original = evt.Content["body"].(string)
						message.Channel = evt.Room.ID
						message.Timestamp = time.Now().String()
						message.Platform = "matrix"

						// Split out command and arguments
						regex := regexp.MustCompile(`^!(\S+)(?: (.+))?$`)
						matches := regex.FindStringSubmatch(message.Message)

						// Extract command
						if len(matches) > 1 && matches[1] != "" {
							message.Command = matches[1]
						}

						// Build arguments list
						if len(matches) > 2 && matches[2] != "" {
							message.Message = strings.TrimSpace(matches[2]) // trim command from message
							message.Arguments = strings.Split(strings.TrimSpace(matches[2]), " ")
						}

						// Loop loaded modules
						var response types.Response
						for _, module := range conf.Modules {
							response = modules.Get(module)(&message)
							handleResponse(response, &message)
						}

						// History module is required
						history.Handle(&message)
					} else {
						fmt.Printf("<%[1]s> %[4]s (%[2]s/%[3]s)\n", evt.Sender, evt.Type, evt.ID, evt.Content["body"])
					}
				default:
					fmt.Println("Unidentified event of type", evt.Type)
				}
			case roomID := <-matrixClient.InviteChan:
				invite := matrixClient.Invites[roomID]
				fmt.Printf("%s invited me to %s (%s)\n", invite.Sender, invite.Name, invite.ID)
				fmt.Println(invite.Accept())
			}
		}
		matrixClient.Stop()
	}()


	// Configure connection
	ircClient = girc.New(girc.Config{
		Server: 		conf.Server,
		Port:   		conf.Port,
		Nick:   		conf.Nick,
		User:   		conf.User,
		// Debug:  		os.Stdout,
		// Out: 		os.Stdout,
		SSL: 			conf.Secure,
		TLSConfig: 		&tls.Config{InsecureSkipVerify: conf.SkipVerify},
		GlobalFormat:	true,
	})

	// Handlers
	ircClient.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		log.Printf("Connected to %s (%s) as nick '%s'", c.Server(), c.ServerVersion(), c.GetNick())

		for _, channel := range conf.Channels {
	    	c.Cmd.Join(channel)
	    	log.Printf("Joined channel %s:%s", c.Server(), channel)
		}
	})


	ircClient.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		message := types.Message{}
		message.Nick = e.Source.Name
		message.Message = e.Trailing
		message.Original = e.Trailing
		message.Timestamp = time.Now().String()
		message.Platform = "irc"

		if len(e.Params) > 0 && girc.IsValidChannel(e.Params[0]) {
			message.Channel = e.Params[0]
		} else {
			message.Private = true
		}

		// Split out command and arguments
		regex := regexp.MustCompile(`^!(\S+)(?: (.+))?$`)
		matches := regex.FindStringSubmatch(message.Message)

		// Extract command
		if len(matches) > 1 && matches[1] != "" {
			message.Command = matches[1]
		}

		// Build arguments list
		if len(matches) > 2 && matches[2] != "" {
			message.Message = strings.TrimSpace(matches[2]) // trim command from message
			message.Arguments = strings.Split(strings.TrimSpace(matches[2]), " ")
		}

		// Loop loaded modules
		var response types.Response
		for _, module := range conf.Modules {
			response = modules.Get(module)(&message)
			handleResponse(response, &message)
		}

		// History module is required
		history.Handle(&message)
	})

	// Connect to server
	if err := ircClient.Connect(); err != nil {
		log.Printf("Error: %s on Server: %s", err, ircClient.Server())
		_, time := ircClient.Uptime()
		log.Printf("%s", time)
	}
}

func handleResponse(response types.Response, original *types.Message) {
	if len(response.Messages) == 0 && response.Message == "" {
		return
	}

	original.Replied = true

	if response.Target == "" {
		if original.Private {
	    	// Reply in PM
	    	response.Target = original.Nick
	    } else {
			// Reply in channel
			response.Target = original.Channel
	    }
	}

	var room *mautrix.Room
	if (original.Platform == "matrix") {
		room = matrixClient.GetRoom(response.Target)
	}

	if len(response.Messages) == 0 {
		response.Messages = append(response.Messages, response.Message)
	}

	for _, message := range response.Messages {
		if (original.Platform == "matrix") {
			if (response.Type == "action") {
				room.Emote(message);
			} else {
				room.Send(message);
			}
		} else {
			if (response.Type == "action") {
				ircClient.Cmd.Action(response.Target, message)
			} else {
				ircClient.Cmd.Message(response.Target, message)
			}
		}
	}
}


func loadConfig() types.Config {

	// Use environment variable 'CONF' or default to './config.json'
	configFile := os.Getenv("CONF")
	if "" == configFile {
		configFile = "config.json"
	}

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Error Loading Configuration (%s): %s", configFile, err)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)

	configuration := types.Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatalf("Error Decoding Configuration (%s): %s", configFile, err)
	}

	// Check Bot Names.
	if !girc.IsValidNick(configuration.Nick) {
		log.Fatalf("Not a valid Nick:", configuration.Nick)
	}

	if !girc.IsValidUser(configuration.User) {
		log.Fatalf("Not a valid User:", configuration.User)
	}

	// Check Admins nicks
	for _, admin := range configuration.Admin {
		if !girc.IsValidNick(admin) {
			log.Fatalf("Not a valid admin Nick:", admin)
		}
	}

	// Check Channels
	for _, channel := range configuration.Channels {
		if !girc.IsValidChannel(channel) {
			log.Fatalf("Not a valid Channel:", channel)
		}
	}

	log.Printf("Configuration file loaded: %s", configFile)

	return configuration
}