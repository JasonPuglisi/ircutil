package ircutil

import (
	"fmt"
	"log"
	"strings"
)

// Log adds a client prefix and logs a message.
func Log(client *Client, message string) {
	log.Printf("%s \t%s\n", GetClientPrefix(client), message)
}

// Logf formats a string, adds a client prefix, and logs a message.
func Logf(client *Client, format string, a ...interface{}) {
	Log(client, fmt.Sprintf(format, a...))
}

// GetClientPrefix formats the client server and user ids into a prefix string
// that's useful terminal output.
func GetClientPrefix(client *Client) string {
	return fmt.Sprintf("[%s/%s]", client.ServerID, client.UserID)
}

// IsChannel determines whether or not a target is a channel. If it is not a
// channel, the target will be a user.
func IsChannel(target string) bool {
	return target[0] == '#'
}

// GetNick isolates a nickname from a source string.
func GetNick(src string) string {
	return src[0:strings.IndexRune(src, '!')]
}
