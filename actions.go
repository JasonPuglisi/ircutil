// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
  "fmt"
  "strings"
)

// Join sends a raw JOIN to attach a client to a channel with an optional
// password. If there is no password, it should be an empty string.
func Join(client *Client, channel string, pass string) {
  sendRaw(client, strings.TrimSpace(fmt.Sprintf("JOIN %s %s", channel, pass)))
}

// Pong sends a raw PONG with a message.
func Pong(client *Client, msg string) {
  sendRawf(client, "PONG %s", msg)
}

// Privmsg sends a raw PRIVMSG to send text to a target.
func Privmsg(client *Client, target string, msg string) {
  sendRawf(client, "PRIVMSG %s :%s", target, msg)
}
