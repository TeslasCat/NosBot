// nosbot.go

package main

import (
	"os"
	"log"
	"strings"
	"time"
	"github.com/lrstanley/girc"
	"crypto/tls"
)

func main() {
	c := loadConfig()

	if c.Debug {
		log.Printf("Printing Configuration file: \n%+v\n", c)
	}

	//  Configure connection
	client := girc.New(girc.Config{
		Server: c.Server,
		Port:   c.Port,
		Nick:   c.Nick,
		User:   c.User,
		Debug:  os.Stdout,
		SSL: 	c.Secure,
		TLSConfig: &tls.Config{InsecureSkipVerify: c.SkipVerify},
	})

	// Handlers
	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
	    c.Cmd.Join("#test")
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		switch {
			case strings.HasPrefix(e.Trailing, "!hello"):
				msgHello(c, e)
		}
	})

	// Connect to server
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

func msgHello (c *girc.Client, e girc.Event) {
	c.Cmd.ReplyTo(e, girc.Fmt("{b}hello{b} {blue}world{c}!"))
	return
}
