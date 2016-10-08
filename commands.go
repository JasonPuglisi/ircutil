package ircutil

import "errors"

// Message stores information about the message the triggered the command.
type Message struct {
	Source  string
	Target  string
	Trigger string
	Args    []string
}

// CmdMap is a map that stores pointers to functions which can be called using
// strings from config.
type CmdMap map[string]CmdFunc

// CmdFunc is a function that's executed for a command, providing necessary
// details to perform an action.
type CmdFunc func(*Client, *Command, *Message)

// InitCommands returns an empty map that can store pointers to functions which
// may be called using strings from config.
func InitCommands() CmdMap {
	return make(CmdMap)
}

// AddCommand adds a function to a command map with a string key.
func AddCommand(cmdMap CmdMap, key string, cmdFunc CmdFunc) {
	cmdMap[key] = cmdFunc
}

// ExecCommand executes a command given a string key, or returns an error if
// the key is not a valid command.
func ExecCommand(client *Client, key string, command *Command,
	message *Message) error {
	if cmdFunc, exists := client.CmdMap[key]; exists {
		go cmdFunc(client, command, message)
	} else {
		return errors.New("executing command: invalid key")
	}
	return nil
}
