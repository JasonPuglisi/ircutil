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

// SendPing sends a ping to a server.
func SendPing(client *Client, msg string) {
	sendRawf(client, "PING :%s", msg)
}

// SendPong replies to a server ping.
func SendPong(client *Client, msg string) {
	sendRawf(client, "PONG :%s", msg)
}

// SendPrivmsg sends a message to a user or channel.
func SendPrivmsg(client *Client, target string, msg string) {
	sendRawf(client, "PRIVMSG %s :%s", target, msg)
}

// SendUser sends initial user details upon server connection.
func SendUser(client *Client, user string, real string) {
	sendRawf(client, "USER %s 0 0 :%s", user, real)
}
