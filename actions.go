// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
  "fmt"
  "strings"
)

// SendJoin attaches to a channel with an optional password. An empty string
// indicates no password.
func SendJoin(client *Client, channel string, pass string) {
  sendRaw(client, strings.TrimSpace(fmt.Sprintf("JOIN %s %s", channel, pass)))
}

// SendNick sets or updates a nickname.
func SendNick(client *Client, nick string) {
  sendRawf(client, "NICK %s", nick)
}

// SendPass authenticates with a password-protected server.
func SendPass(client *Client, pass string) {
  sendRawf(client, "PASS :%s", pass)
}

// SendPong replies to a server ping.
func SendPong(client *Client, msg string) {
  sendRawf(client, "PONG %s", msg)
}

// SendPrivmsg sends a message to a user or channel.
func SendPrivmsg(client *Client, target string, msg string) {
  sendRawf(client, "PRIVMSG %s :%s", target, msg)
}

// SendUser sends initial user details upon server connection. It parses a
// string for the two possible initial mode characters 'w' and 'i'. No other
// characters may be present.
func SendUser(client *Client, user string, mode string, real string) {
  intMode := 0
  if strings.IndexRune(mode, 'w') < 0 {
    intMode += 4
  }
  if strings.IndexRune(mode, 'i') < 0 {
    intMode += 8
  }

  sendRawf(client, "USER %s %d * :%s", user, mode, real)
}
