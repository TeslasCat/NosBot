package main

import (
	"os"
	"log"
	"strings"
	"encoding/json"
	"crypto/tls"
	"github.com/lrstanley/girc"

	"./modules/notes"
)

func main() {

	conf := LoadConfig()

	if conf.Debug {
		log.Printf("Printing Configuration file: \n%+v\n", conf)
	}

	//  Configure connection
	client := girc.New(girc.Config{
		Server: conf.Server,
		Port:   conf.Port,
		Nick:   conf.Nick,
		User:   conf.User,
		// Debug:  os.Stdout,
		// Out: 	os.Stdout,
		SSL: 	conf.Secure,
		TLSConfig: &tls.Config{InsecureSkipVerify: conf.SkipVerify},
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
		switch {
			case strings.HasPrefix(e.Trailing, "!note"):
				notes.New(c, e)
			case strings.HasPrefix(e.Trailing, "!list"):
				notes.List(c, e)
		}
	})

	// Connect to server
	if err := client.Connect(); err != nil {
		log.Printf("Error: %s on Server: %s", err, client.Server())
		_, time := client.Uptime()
		log.Printf("%s", time)
	}
}

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
}

func LoadConfig() Config {

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

	configuration := Config{}
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