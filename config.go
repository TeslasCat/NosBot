// config.go

package main

import (
    "os"
    "log"
    "encoding/json"
    "github.com/lrstanley/girc"
)

// https://mholt.github.io/json-to-go/ <3
type Config struct {
	Server         string   `json:"server"`
	Channels       []string `json:"channels"`
	Nick           string   `json:"nick"`
	User           string   `json:"user"`
	Nickserv       string   `json:"nickserv"`
	Debug          bool     `json:"debug"`
	Port           int      `json:"port"`
	Secure         bool     `json:"secure"`
	Admin          []string `json:"admin"`
	WordnikAPI     string   `json:"wordnik_api"`
	Greeting       []string `json:"greeting"`
	GreetingIgnore []string `json:"greeting-ignore"`
}

func loadConfig() Config {

	// Use environment variable 'CONF' or default to './config.json'
	configFile := os.Getenv("CONF")
	if "" == configFile {
		configFile = "config.json"
	}

	file, _ := os.Open(configFile)
	defer file.Close()

	decoder := json.NewDecoder(file)

	configuration := Config{}
	err := decoder.Decode(&configuration)

	if err != nil {
		log.Printf("Error loading configuration (%s): %s", configFile, err)
		os.Exit(1)
	}

	// Check Bot Names.
	if !girc.IsValidNick(configuration.Nick) {
		log.Printf("Not a valid Nick:", configuration.Nick)
		os.Exit(1)
	}

	if !girc.IsValidNick(configuration.User) {
		log.Printf("Not a valid User:", configuration.User)
		os.Exit(1)
	}

	// Check Admins nicks
	for _, admins := range configuration.Admin {
		if !girc.IsValidNick(admins) {
			log.Printf("Not a valid admin Nick:", admins)
			os.Exit(1)
		}
	}

	// Check Channels
	for _, channels := range configuration.Channels {
		if !girc.IsValidChannel(channels) {
			log.Printf("Not a valid Channel:", channels)
			os.Exit(1)
		}
	}

	log.Printf("Configuration file loaded.", file.Name)

	return configuration
}