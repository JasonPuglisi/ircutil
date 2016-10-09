package ircutil

import (
	"fmt"
	"log"
	"strings"
)

// Log adds a client prefix and logs a message.
func Log(client *Client, message string) {
	log.Printf("%s \t%s\n", getClientPrefix(client), message)
}

// Logf formats a string, adds a client prefix, and logs a message.
func Logf(client *Client, format string, a ...interface{}) {
	Log(client, fmt.Sprintf(format, a...))
}

// getClientPrefix formats the client server and user ids into a prefix string
// that's useful terminal output.
func getClientPrefix(client *Client) string {
	return fmt.Sprintf("[%s/%s]", client.ServerID, client.UserID)
}

// isChannel determines whether or not a target is a channel. If it is not a
// channel, the target will be a user.
func isChannel(target string) bool {
	return target[0] == '#'
}

// getNick isolates a nickname from a source string.
func getNick(src string) string {
	return src[0:strings.IndexRune(src, '!')]
}
