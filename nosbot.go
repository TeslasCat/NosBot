// nosbot.go

package main

import (
	"os"
	"log"
	"strings"
	"time"
	"github.com/lrstanley/girc"
)

func main() {
	c := loadConfig()

	if c.Debug {
		log.Printf("Printing Configuration file: \n%+v\n", c)
	}

	os.Exit(3)

	client := girc.New(girc.Config{
		Server: c.Server,
		Port:   c.Port,
		Nick:   c.Nick,
		User:   c.User,
		Debug:  os.Stdout,
		SSL: 	c.Secure,
	})

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
	    c.Cmd.Join("#test")
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
	    if strings.HasPrefix(e.Trailing, "!hello") {
	        c.Cmd.ReplyTo(e, girc.Fmt("{b}hello{b} {blue}world{c}!"))
	        return
	    }
	})

	for {
	    if err := client.Connect(); err != nil {
	        log.Printf("error: %s", err)

	        log.Println("reconnecting in 30 seconds...")
	        time.Sleep(30 * time.Second)
	    } else {
	        return
	    }
	}
}
