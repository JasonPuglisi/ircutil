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
  msg = strings.TrimSuffix(msg, "\r\n")
  if (client.Debug) {
    fmt.Printf(">> %s\n", msg)
  }

  // Set empty source and split message into tokens. Update source and remove
  // it from tokens if found.
  src, tokens := "", strings.Split(msg, " ")
  if (tokens[0][0] == ':') {
    src, tokens = tokens[0], tokens[1:]
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
  }
}

// handleCommand takes a command to determine the correct action to take after
// receiving a message from a server.
func handleCommand(client *Client, src string, cmd string, tokens []string) {
  switch cmd {
  // PING is sent by servers upon connection and at regular intervals. We will
  // send the same string back.
  case "PING":
    Pong(client, strings.Join(tokens, " "))
  }
}
