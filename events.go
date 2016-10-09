// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
	"fmt"
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
		Logf(client, "<= %s", msg)
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
	// NICK is sent when a nickname is updated. Update client's state if it
	// belongs to the client.
	case "NICK":
		if client.Nick == getNick(src) {
			client.Nick = strings.Join(tokens, " ")[1:]
		}
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
			trigger := fmt.Sprintf("%s%s", s.Symbol, t)
			if validateCommand(client, s, trigger, src, target, cmd) {
				// Make sure command call has enough arguments, or error if not.
				if checkArgs(c.Arguments, tokens) {
					// Execute command that was found, or error if the function key is not
					// valid.
					err := ExecCommand(client, c.Function, c, &Message{src, target,
						trigger, tokens})
					if err != nil {
						Log(client, err.Error())
					}
				} else {
					SendResponse(client, src, target,
						fmt.Sprintf("Invalid arguments. Usage: %s %s", trigger,
							c.Arguments))
				}
			}
		}
	}
}

// validateCommand makes sure the user has the trigger matches the command,
// the scope matches the command, and the user admin permissions match the
// command. It does not check arguments, leaving that to be done separately.
func validateCommand(client *Client, settings *Settings, trigger string,
	src string, target string, cmd string) bool {
	// Ensure the trigger matches the command, taking into account the
	// case-sensitivity setting.
	trigger = processTrigger(client, trigger)
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

	// Ensure user admin permissions match the command.
	if settings.Admin {
		adminMatch := false
		nick := getNick(src)
		for i := range client.Admins {
			a := &client.Admins[i]
			if *a == nick {
				adminMatch = true
			}
		}
		if !adminMatch {
			return false
		}
	}

	// All matches have passed, so return true.
	return true
}

// processTrigger fills in the nickname variable if it exists within a trigger.
func processTrigger(client *Client, trigger string) string {
	return strings.Replace(trigger, "%NICK%", client.Nick, 1)
}

// checkArgs determines whether or not a command called by a user has enough
// arguments to be executed.
func checkArgs(list string, args []string) bool {
	// Increment number of needed arguments based on chevron-enclosed tokens in
	// the list.
	needed := 0
	if len(list) > 0 {
		for _, arg := range strings.Split(list, " ") {
			if arg[0] == '<' && arg[len(arg)-1] == '>' {
				needed++
			}
		}
	}
	return len(args) >= needed
}
