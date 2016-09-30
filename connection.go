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
  "strings"
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
  nick  string
  user  string
  real  string
  modes string
}

// Client stores initial server/user details, client status, client channels,
// and active client connection.
type Client struct {
  Active  bool
  Debug   bool
  Ready   func(*Client)
  server  *Server
  user    *User
  conn    net.Conn
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
// Modes must be a string containing only characters 'w' and 'i' or neither.
func CreateUser(nick string, uname string, real string, modes string) (*User,
    error) {
  // Error if nickname is empty.
  if len(nick) < 1 {
    return &User{}, errors.New("creating user: nickname too short")
  }

  // Error if modes are invalid.
  runes, posW, posI := len(modes), strings.IndexRune(modes, 'w'),
    strings.IndexRune(modes, 'i')
  if (runes == 1 && posW == posI) || (runes == 2 && (posW < 0 || posI < 0)) ||
      runes > 2 {
    return &User{}, errors.New("creating user: mode string invalid")
  }

  return &User{nick, uname, real, modes}, nil
}

// EstablishConnection establishes a connection to the specified IRC server
// using the specified user information. It sends initial messages as required
// by the IRC protocol.
func EstablishConnection(server *Server, user *User, ready func(*Client),
    debug bool) (*Client, error) {
  // Attempt connection establishment. Use TLS if secure is specified. Timeout
  // after one minute.
  var conn net.Conn; var err error
  if server.secure {
    conn, err = tls.DialWithDialer(&(net.Dialer{Timeout: time.Duration(1) *
      time.Minute}), "tcp", fmt.Sprintf("%s:%d", server.host, server.port),
      nil)
  } else {
    conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", server.host,
      server.port), time.Duration(1) * time.Minute)
  }
  if err != nil {
    return &Client{}, err
  }
  fmt.Printf("Connected to server \"%s:%d\" (%s)\n", server.host, server.port,
    conn.RemoteAddr())

  // Create client with with server, user, connection, intitialization
  // function, and debug setting. Start reading from server.
  client := Client{true, debug, ready, server, user, conn}
  go readLoop(&client)
  go pingLoop(&client)

  // Send required user registration messages to server, including password if
  // specified.
  if len(server.pass) > 0 {
    SendPass(&client, server.pass)
  }
  SendNick(&client, user.nick)
  SendUser(&client, user.user, user.modes, user.real)

  return &client, nil;
}

// readLoop is a goroutine used to read data from the server a client is
// connected to.
func readLoop(client *Client) {
  // Create reader for server data, and loop reads until client is stopped.
  reader := bufio.NewReader(client.conn)
  for client.Active {
    msg, _ := reader.ReadString('\n')

    parseMessage(client, msg)
  }
}

// pingLoop is a goroutine used to send periodic pings to the server a client
// is connected to in order to keep that connection alive.
func pingLoop(client *Client) {
  ticker := time.NewTicker(time.Minute)
  for range ticker.C {
    SendPing(client, time.Now().String())
  }
}

// sendRawf formats a string to create a raw IRC message that is sent to the
// client's server. It appends necessary line endings.
func sendRawf(client *Client, format string, a ...interface{}) {
  fmt.Fprintf(client.conn, format + "\r\n", a...)
  if client.Debug {
    log.Printf(">> " + format + "\n", a...)
  }
}

// sendRaw sends a raw IRC message to the client's server. It appends necessary
// line endings.
func sendRaw(client *Client, msg string) {
  fmt.Fprint(client.conn, msg + "\r\n")
  if client.Debug {
    log.Println(">>", msg)
  }
}
