// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
	"math/rand"
	"strconv"
)

// SendJoin attaches to a channel with an optional password. An empty string
// indicates no password.
func SendJoin(client *Client, channel string, pass string) {
	if pass != "" {
		pass = " " + pass
	}
	sendRawf(client, "JOIN %s%s", channel, pass)
}

// SendModeUser updates user modes.
func SendModeUser(client *Client, modes string) {
	sendRawf(client, "MODE %s %s", client.Nick, modes)
}

// SendNick sets or updates a nickname.
func SendNick(client *Client, nick string) {
	client.Nick = nick
	sendRawf(client, "NICK %s", nick)
}

// SendNickRandom sets or updates a nickname to a random one.
func SendNickRandom(client *Client) {
	SendNick(client, "Inami"+strconv.Itoa(rand.Intn(10000)))
}

// SendNickservPass authenticates a nickname with Nickserv.
func SendNickservPass(client *Client, pass string) {
	SendPrivmsg(client, "nickserv", "identify "+pass)
}

// SendNotice sends a notice to a user or channel.
func SendNotice(client *Client, target string, msg string) {
	sendRawf(client, "NOTICE %s :%s", target, msg)
}

// SendPart detaches from a channel with an optional message. An empty string
// indicates no message.
func SendPart(client *Client, channel string, msg string) {
	if msg != "" {
		msg = " :" + msg
	}
	sendRawf(client, "PART %s%s", channel, msg)
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

// SendResponse determines whether a message should be sent to a user or
// or channel, and sends the message accordingly.
func SendResponse(client *Client, src string, target string, msg string) {
	if isChannel(target) {
		SendPrivmsg(client, target, msg)
	} else {
		SendPrivmsg(client, getNick(src), msg)
	}
}

// SendUser sends initial user details upon server connection.
func SendUser(client *Client, user string, real string) {
	sendRawf(client, "USER %s 0 0 :%s", user, real)
}
