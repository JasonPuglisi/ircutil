// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

// Settings stores command settings that can be used as defaults or overrode on
// a per-command basis.
type Settings struct {
	// (Optional) Whether or not a command is case-sensitive. Default: true
	CaseSensitive bool `json:"caseSensitive"`
	// (Optional) Symbol/string that must be prefixed to a command trigger for
	// it to be detected. Default: "!"
	Symbol string `json:"symbol"`
	// (Optional) Array of places a command can be used. Contains "channel" for
	// in a channel, and/or "direct" for in a private message to the client.
	// Default: ["channel"]
	Scope []string `json:"scope"`
	// (Optional) Whether or not a command can only be used by client admins.
	// Default: false
	Admin bool `json:"admin"`
}

// Server stores server connection settings.
type Server struct {
	// String ID of the server. Used to identify the server in a client. Must be
	// alphanumeric.
	ID string `json:"id"`
	// Hostname of server. Can be an IP address or domain name.
	Host string `json:"host"`
	// (Optional) Port of server. Default: 6697
	Port uint16 `json:"port"`
	// (Optional) Whether or not the server port should be connected to using
	// SSL/TLS. Default: true if port is 6697, false otherwise
	Secure bool `json:"secure"`
}

// User stores user settings.
type User struct {
	// String ID of user. Used to identify the user in a client. Must be
	// alphanumeric.
	ID string `json:"id"`
	// Nickname of client.
	Nick string `json:"nick"`
	// (Optional) Username of client. Default: Same as nickname lowercased
	User string `json:"user"`
	// (Optional) Realname of client. Default: Same as nickname
	Real string `json:"real"`
}

// Client stores client/connection settings and credentials.
type Client struct {
	// ID of server to use in connection.
	ServerID string `json:"serverId"`
	// ID of user to use in connection.
	UserID string `json:"userId"`
	// (Optional) List of channels to join upon connection and authentication
	// (if specified). Channels must be prefixed with "#" (e.g., "#channel").
	// Channels with a password must have a space between the channel name and
	// password (e.g., "#channel pass"). Default: []
	Channels []string `json:"channels"`
	// (Optional) String of user modes to be set upon connection and
	// authentication (if specified). Must have a "+" before all modes to be
	// set and a "-" before all modes to be unset (e.g., "+i-x"). Default: ""
	Modes string `json:"modes"`
	// (Optional) List of client admin nicknames able to run commands set to
	// admin-only. Default: []
	Admins []string `json:"admins"`
	// (Optional) Authentication credentials to used in connection.
	// Defaults: All nested defaults
	Authentication `json:"authentication"`
	// Other non-configurable values.
	CmdMap
	Commands []Command
	Debug    bool
	Ready    func(*Client)
	Done     chan bool
	Server   *Server
	User     *User
	Conn     net.Conn
	Nick     string
}

// Authentication stores authentication credentials for servers and nicknames.
type Authentication struct {
	// (Optional) Server password to connect to a server. Empty string for
	// none. Default: ""
	ServerPassword string `json:"serverPassword"`
	// (Optional) Nickserv password to identify user with nickserv. Empty
	// string for none. Default: ""
	Nickserv string `json:"nickserv"`
}

// Command stores command triggers, execution details, and settings, with the
// ability to override default settings.
type Command struct {
	// List of strings that will trigger the command. Triggers must not contain
	// the command symbol, as it will be checked for automatically.
	Triggers []string `json:"triggers"`
	// Function that will be executed when the command is triggered.
	Function string `json:"function"`
	// (Optional) String of arguments that must follow the command. Mandatory
	// arguments must have chevrons around them, and optional arguments must have
	// square brackets around them (e.g., "<arg1> <arg2> [arg3]"). Default: ""
	Arguments string `json:"arguments"`
	// (Optional) Command settings to override default command settings. See
	// default command settings (above) for descriptions of each setting.
	// Defaults: Default command settings
	Settings `json:"settings"`
}

// EstablishConnection establishes a connection to the specified IRC server
// using the specified user information. It sends initial messages as required
// by the IRC protocol.
func EstablishConnection(client *Client) error {
	// Error if server or user id is empty or non-alphanumeric.
	r, _ := regexp.Compile("^[0-9A-Za-z]+$")
	matched := r.MatchString(client.ServerID)
	if !matched {
		return errors.New("establishing connection: server id invalid")
	}
	matched = r.MatchString(client.UserID)
	if !matched {
		return errors.New("establishing connection: user id invalid")
	}

	// Error if hostname is empty.
	if len(client.Server.Host) < 1 {
		return errors.New("establishing connection: hostname too short")
	}

	// Error if port is too low.
	if client.Server.Port == 0 {
		return errors.New("establishing connection: port too low")
	}

	// Error if nickname is empty.
	if len(client.User.Nick) < 1 {
		return errors.New("establishing connection: nickname too short")
	}

	// Attempt connection establishment. Use TLS if secure is specified. Timeout
	// after one minute.
	var conn net.Conn
	var err error
	if client.Server.Secure {
		conn, err = tls.DialWithDialer(&(net.Dialer{Timeout: time.Minute}), "tcp",
			fmt.Sprintf("%s:%d", client.Server.Host, client.Server.Port), nil)
	} else {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d",
			client.Server.Host, client.Server.Port), time.Minute)
	}
	if err != nil {
		return err
	}
	Logf(client, "Connected to server %s:%d (%s)", client.Server.Host,
		client.Server.Port, conn.RemoteAddr())

	// Update connection in client and start reading from server and pinging
	// periodically.
	client.Conn = conn
	go readLoop(client)
	go pingLoop(client)

	// Send required user registration messages to server, including password if
	// specified.
	if len(client.Authentication.ServerPassword) > 0 {
		SendPass(client, client.Authentication.ServerPassword)
	}
	SendNick(client, client.User.Nick)
	SendUser(client, client.User.User, client.User.Real)

	return nil
}

// readLoop is a goroutine used to read data from the server a client is
// connected to.
func readLoop(client *Client) {
	// Create reader for server data.
	reader := bufio.NewReader(client.Conn)
	for {
		select {
		// If the client is done, stop the goroutine.
		case <-client.Done:
			return
		// Loop reads from server and set client to done if a timeout or error is
		// encountered.
		default:
			// Make sure a message is received in at most three minutes, since we
			// send a ping every two.
			client.Conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
			msg, err := reader.ReadString('\n')
			client.Conn.SetReadDeadline(time.Time{})
			if err != nil {
				Log(client, err.Error())
				close(client.Done)
			} else {
				parseMessage(client, msg)
			}
		}
	}
}

// pingLoop is a goroutine used to send periodic pings to the server a client
// is connected to in order to keep that connection alive.
func pingLoop(client *Client) {
	// Create ticker to send pings every two minutes.
	ticker := time.NewTicker(time.Minute * 2)
	for {
		select {
		// If the client is done, stop the time and goroutine.
		case <-client.Done:
			ticker.Stop()
			return
		// Loop pings to keep connection alive.
		case <-ticker.C:
			SendPing(client, strconv.FormatInt(time.Now().UnixNano(), 10))
		}
	}
}
