// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// Server stores server connection details.
type Server struct {
	host   string
	port   uint16
	secure bool
	pass   string
}

// User stores initial connection user details.
type User struct {
	nick string
	user string
	real string
}

// Client stores initial server/user details, client status, client channels,
// and active client connection.
type Client struct {
	Debug  bool
	Ready  func(*Client)
	Done   chan bool
	server *Server
	user   *User
	conn   net.Conn
}

// CreateServer creates and returns a server for use in connections.
// (Host)name is an IP address or FQDN. IPv6 addresses must be in [brackets].
// Secure is true if the specified port uses TLS/SSL.
// (Pass)word must be an empty string if there is no server password.
func CreateServer(host string, port uint16, secure bool, pass string) (*Server,
	error) {
	// Error if hostname is empty.
	if len(host) < 1 {
		return &Server{}, errors.New("creating server: hostname too short")
	}

	return &Server{host, port, secure, pass}, nil
}

// CreateUser creates and returns a user for use in connections.
// Username (uname) and (real)name will default to (nick)name if empty strings.
func CreateUser(nick string, uname string, real string) (*User, error) {
	// Error if nickname is empty.
	if len(nick) < 1 {
		return &User{}, errors.New("creating user: nickname too short")
	}

	// Set username and realname to nickname if empty.
	if len(uname) < 1 {
		uname = nick
	}
	if len(real) < 1 {
		real = nick
	}

	return &User{nick, uname, real}, nil
}

// EstablishConnection establishes a connection to the specified IRC server
// using the specified user information. It sends initial messages as required
// by the IRC protocol.
func EstablishConnection(server *Server, user *User, ready func(*Client),
	debug bool) (*Client, error) {
	// Attempt connection establishment. Use TLS if secure is specified. Timeout
	// after one minute.
	var conn net.Conn
	var err error
	if server.secure {
		conn, err = tls.DialWithDialer(&(net.Dialer{Timeout: time.Minute}), "tcp",
			fmt.Sprintf("%s:%d", server.host, server.port), nil)
	} else {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", server.host,
			server.port), time.Minute)
	}
	if err != nil {
		return &Client{}, err
	}
	log.Printf("Connected to server \"%s:%d\" (%s)\n", server.host, server.port,
		conn.RemoteAddr())

	// Create client with with server, user, connection, intitialization
	// function, and debug setting. Start reading from server.
	client := Client{debug, ready, make(chan bool, 1), server, user, conn}
	go readLoop(&client)
	go pingLoop(&client)

	// Send required user registration messages to server, including password if
	// specified.
	if len(server.pass) > 0 {
		SendPass(&client, server.pass)
	}
	SendNick(&client, user.nick)
	SendUser(&client, user.user, user.real)

	return &client, nil
}

// readLoop is a goroutine used to read data from the server a client is
// connected to.
func readLoop(client *Client) {
	// Create reader for server data.
	reader := bufio.NewReader(client.conn)
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
			client.conn.SetReadDeadline(time.Now().Add(time.Minute * 3))
			msg, err := reader.ReadString('\n')
			client.conn.SetReadDeadline(time.Time{})
			if err != nil {
				log.Println(err)
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

// sendRawf formats a string to create a raw IRC message that is sent to the
// client's server. It appends necessary line endings.
func sendRawf(client *Client, format string, a ...interface{}) {
	fmt.Fprintf(client.conn, format+"\r\n", a...)
	if client.Debug {
		log.Printf(">> "+format+"\n", a...)
	}
}

// sendRaw sends a raw IRC message to the client's server. It appends necessary
// line endings.
func sendRaw(client *Client, msg string) {
	fmt.Fprint(client.conn, msg+"\r\n")
	if client.Debug {
		log.Println(">>", msg)
	}
}
