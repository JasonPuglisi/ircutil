// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// parseMessage extracts the sender from a raw IRC message and determines if
// the message is sending a command or response code. It also extracts the
// response code if applicable.
func parseMessage(client *Client, msg string) {
	// Remove line ending and print message to console for debugging.
	msg = strings.TrimSpace(strings.TrimSuffix(msg, "\r\n"))
	if client.Debug {
		log.Printf("%s \t<= %s\n", GetClientPrefix(client), msg)
	}

	// Set empty source and split message into tokens. Update source and remove
	// it from tokens if found.
	src, tokens := "", strings.Split(msg, " ")
	if tokens[0][0] == ':' {
		src, tokens = tokens[0][1:], tokens[1:]
	}

	// Attempt to parse first token as number. If successful, handle the message
	// as a response. If not, handle the message as a command.
	if code, err := strconv.Atoi(tokens[0]); err == nil {
		handleResponse(client, src, code, tokens[1:])
	} else {
		handleCommand(client, src, tokens[0], tokens[1:])
	}
}

// handleResponse takes a response code to determine the correct action to take
// after receiving a message from a server.
func handleResponse(client *Client, src string, code int, tokens []string) {
	switch code {
	// 004 RPL_MYINFO is the last mandatory message to be sent after a client
	// registers with a server, meaning we can now start performing actions.
	case 4:
		client.Ready(client)
	// 433 ERR_NICKNAMEINUSE is send when the client tries to change its nick
	// to one another user using, forcing us to choose a random one.
	case 433:
		SendNickRandom(client)
	}
}

// handleCommand takes a command to determine the correct action to take after
// receiving a message from a server.
func handleCommand(client *Client, src string, cmd string, tokens []string) {
	switch cmd {
	// NICK is sent by servers when they force update a client's nickname,
	// leaving the client to update its internal state.
	case "NICK":
		client.Nick = strings.Join(tokens, " ")[1:]
	// PING is sent by servers upon connection and at regular intervals. We will
	// send the same string back.
	case "PING":
		SendPong(client, strings.Join(tokens, " ")[1:])
	case "PRIVMSG":
		handleMessage(client, src, tokens[0], tokens[1][1:], tokens[2:])
	}
}

// handleMessage checks a message for a valid command, end executes the command
// if found.
func handleMessage(client *Client, src string, target string, cmd string,
	tokens []string) {
	// Loop through all commands.
	for i := range client.Commands {
		c := &client.Commands[i]
		s := &c.Settings
		// Loop through all triggers for each command.
		for j := range c.Triggers {
			t := c.Triggers[j]
			// Prepend symbol to trigger, and try to match it to the user message.
			// Ignore cases if case sensitivity is off. Validate channel/user scope.
			trigger := s.Symbol + t
			if validateCommand(t, target, cmd, s) {
				// Make sure command call has enough arguments, or error if not.
				if checkArgs(c.Arguments, tokens) {
					// Execute command that was found, or error if the function key is not
					// valid.
					err := ExecCommand(client, c.Function, c, &Message{src, target,
						trigger, tokens})
					if err != nil {
						log.Println(err)
					}
				} else {
					SendResponse(client, src, target, fmt.Sprintf("Invalid arguments. "+
						"Usage: %s %s", trigger, c.Arguments))
				}
			}
		}
	}
}

// validateCommand makes sure the user has the trigger matches the command,
// the scope matches the command, and the user admin permissions match the
// command. It does not check arguments, leaving that to be done separately.
func validateCommand(trigger string, target string, cmd string,
	settings *Settings) bool {
	// Ensure the trigger matches the command, taking into account the
	// case-sensitivity setting.
	triggerMatch := trigger == cmd || (!settings.CaseSensitive &&
		strings.ToLower(trigger) == strings.ToLower(cmd))
	if !triggerMatch {
		return false
	}

	// Ensure the scope matches the command.
	scopeMatch := false
	for i := range settings.Scope {
		s := &settings.Scope[i]
		if (*s == "channel" && isChannel(target)) || (*s == "direct" &&
			!isChannel(target)) {
			scopeMatch = true
		}
	}
	if !scopeMatch {
		return false
	}

	// All matches have passed, so return true.
	return true
}

// checkArgs determines whether or not a command called by a user has enough
// arguments to be executed.
func checkArgs(list string, args []string) bool {
	// Increment number of needed arguments based on chevron-enclosed tokens in
	// the list.
	needed := 0
	for _, arg := range strings.Split(list, " ") {
		if arg[0] == '<' && arg[len(arg)-1] == '>' {
			needed++
		}
	}
	return len(args) >= needed
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
