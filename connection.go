// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
  "bufio"
  "crypto/tls"
  "errors"
  "fmt"
  "net"
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
  mode byte
}

// Client stores initial server/user details, client status, client channels,
// and active client connection.
type Client struct {
  Active  bool
  server  Server
  user    User
  conn    net.Conn
  errCh   chan error
}

// CreateServer creates and returns a server for use in connections.
// (Host)name is an IP address or FQDN. IPv6 addresses must be in [brackets].
// Secure is true if the specified port uses TLS/SSL.
// (Pass)word must be an empty string if there is no server password.
func CreateServer(host string, port uint16, secure bool, pass string) (Server,
    error) {
  // Error if hostname is empty.
  if len(host) < 1 {
    return Server{}, errors.New("creating server: hostname too short")
  }

  return Server{host, port, secure, pass}, nil
}

// CreateUser creates and returns a user for use in connections.
// Username (uname) and (real)name will default to (nick)name if empty strings.
// Mode must be 0 for none, 4 for +w, 8 for +i, or 12 for +iw.
func CreateUser(nick string, uname string, real string, mode byte) (User,
    error) {
  // Error if nickname is empty.
  if len(nick) < 1 {
    return User{}, errors.New("creating user: nickname too short")
  }

  // Error if bitmaskk is not 0, 4, 8, or 12.
  if mode != 0 && mode != 4 && mode != 8 && mode != 12 {
    return User{}, errors.New("creating user: mode bitmask invalid")
  }

  return User{nick, uname, real, mode}, nil
}

// EstablishConnection establishes a connection to the specified IRC server
// using the specified user information. It sends initial messages as required
// by the IRC protocol.
func EstablishConnection(server Server, user User) (Client, error) {
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
    return Client{}, err
  }
  fmt.Printf("Connected to server \"%s\" (%s)\n", server.host,
    conn.RemoteAddr())

  // Create client with with server, user, and connection. Start reading from
  // server.
  client := Client{true, server, user, conn, make(chan error)}
  go readLoop(client)

  // Send required user registration messages to server, including password if
  // specified.
  if len(server.pass) > 0 {
    sendRawf(client, "PASS :%s", server.pass)
  }
  sendRawf(client, "NICK %s", user.nick)
  sendRawf(client, "USER %s %d * :%s", user.user, user.mode, user.real)

  return client, nil;
}


// readLoop is a goroutine used to read data from the server a client is
// connected to.
func readLoop(client Client) {
  // Create reader for server data, and loop reads until client is stopped.
  reader := bufio.NewReader(client.conn)
  for client.Active {
    msg, err := reader.ReadString('\n')
    if err != nil {
      client.errCh <- err
    }

    fmt.Printf(msg)
  }
}

// sendRawf formats a string to create a raw IRC message that is sent to the
// client's server. It appends necessary line endings.
func sendRawf(client Client, format string, a ...interface{}) {
  fmt.Fprintf(client.conn, format + "\r\n", a...);
}

// sendRaw sends a raw IRC message to the client's server. It appends necessary
// line endings.
func sendRaw(client Client, msg string) {
  fmt.Fprint(client.conn, msg + "\r\n");
}
