// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
	"fmt"
	"math/rand"
)

// SendJoin attaches to a channel with an optional password. An empty string
// indicates no password.
func SendJoin(client *Client, channel string, pass string) {
	if len(pass) > 0 {
		pass = fmt.Sprintf(" %s", pass)
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
	SendNick(client, fmt.Sprintf("Inami%d", rand.Intn(100000)))
}

// SendNickservPass authenticates a nickname with Nickserv.
func SendNickservPass(client *Client, pass string) {
	SendPrivmsg(client, "nickserv", fmt.Sprintf("identify %s", pass))
}

// SendNotice sends a notice to a user or channel.
func SendNotice(client *Client, target string, msg string) {
	sendRawf(client, "NOTICE %s :%s", target, msg)
}

// SendPart detaches from a channel with an optional message. An empty string
// indicates no message.
func SendPart(client *Client, channel string, msg string) {
	if len(msg) > 0 {
		msg = fmt.Sprintf(" :%s", msg)
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
	if IsChannel(target) {
		SendPrivmsg(client, target, msg)
	} else {
		SendPrivmsg(client, GetNick(src), msg)
	}
}

// SendUser sends initial user details upon server connection.
func SendUser(client *Client, user string, real string) {
	sendRawf(client, "USER %s 0 0 :%s", user, real)
}

// sendRaw sends a raw IRC message to the client's server. It appends necessary
// line endings.
func sendRaw(client *Client, msg string) {
	fmt.Fprintf(client.Conn, "%s\r\n", msg)
	if client.Debug {
		Logf(client, "=> %s", msg)
	}
}

// sendRawf formats a string to create a raw IRC message that is sent to the
// client's server. It appends necessary line endings.
func sendRawf(client *Client, format string, a ...interface{}) {
	sendRaw(client, fmt.Sprintf(format, a...))
}
