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
	SkipVerify     bool     `json:"skipVerify"`
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

	if !girc.IsValidNick(configuration.User) {
		log.Fatalf("Not a valid User:", configuration.User)
	}

	// Check Admins nicks
	for _, admins := range configuration.Admin {
		if !girc.IsValidNick(admins) {
			log.Fatalf("Not a valid admin Nick:", admins)
		}
	}

	// Check Channels
	for _, channels := range configuration.Channels {
		if !girc.IsValidChannel(channels) {
			log.Fatalf("Not a valid Channel:", channels)
		}
	}

	log.Printf("Configuration file loaded.", file.Name)

	return configuration
}
