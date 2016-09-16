// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
  "fmt"
  "strings"
)

// Pong sends a raw PONG to a server with the specified message.
func Pong(client *Client, msg string) {
  sendRawf(client, "PONG %s", msg)
}

// Join attaches a client to a channel with an optional password.
func Join(client *Client, channel string, pass string) {
  sendRaw(client, strings.TrimSpace(fmt.Sprintf("JOIN %s %s", channel, pass)))
}
