package main

import (
	"os"
	"log"
	"encoding/json"
	"crypto/tls"
	"github.com/lrstanley/girc"

	"./modules/notes"
	"./types"
)

var client *girc.Client

func main() {

	conf := loadConfig()

	if conf.Debug {
		log.Printf("Printing Configuration file: \n%+v\n", conf)
	}

	// Configure connection
	client = girc.New(girc.Config{
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
	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		log.Printf("Connected to %s (%s) as nick '%s'", c.Server(), c.ServerVersion(), c.GetNick())

		for _, channel := range conf.Channels {
	    	c.Cmd.Join(channel)
	    	log.Printf("Joined channel %s:%s", c.Server(), channel)
		}
	})


	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		message := types.Message{e.Source.Name, "", e.Trailing, e.Timestamp.String(), false}
		if len(e.Params) > 0 && girc.IsValidChannel(e.Params[0]) {
			message.Channel = e.Params[0]
		} else {
			message.Private = true
		}

		// Loop loaded modules
		var response types.Response
		for _, module := range conf.Modules {
			if module == "notes" {
				response = notes.Handle(message)
				handleResponse(response, message);
			}
		}
	})

	// Connect to server
	if err := client.Connect(); err != nil {
		log.Printf("Error: %s on Server: %s", err, client.Server())
		_, time := client.Uptime()
		log.Printf("%s", time)
	}
}

func handleResponse(response types.Response, original types.Message) {
	if len(response.Messages) == 0 {
		return
	}

	if response.Type == "" {
		response.Type = "message"
	}

	if response.Target == "" {
		if original.Private {
	    	// Reply in PM
	    	response.Target = original.Nick
	    } else {
			// Reply in channel
			response.Target = original.Channel
	    }
	}

	if (response.Type == "action") {
		for _, message := range response.Messages {
			client.Cmd.Action(response.Target, message)
		}
	} else {
		for _, message := range response.Messages {
			client.Cmd.Message(response.Target, message)
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