// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
  "bufio"
  "errors"
  "fmt"
  "net"
)

// Server stores provided server information.
type Server struct {
  Host   string
  Port   uint16
  Secure bool
  Pass   string
}

// User stores provided user information, and is not always current.
type User struct {
  Nick string
  User string
  Real string
  Mode byte
}

// Client stores client/connection information, including original Server and
// User.
type Client struct {
  Active bool
  Server Server
  User   User
  Conn   net.Conn
}

// CreateServer creates and returns a server for use in connections.
// (Host)name is an IP address or FQDN. IPv6 addresses must be in [brackets].
// Secure is true if the specified port uses TLS/SSL.
// (Pass)word must be an empty string if there is no server password.
func CreateServer(host string, port uint16, secure bool, pass string) (Server,
    error) {
  // Error if hostname is empty
  if len(host) < 1 {
    return Server{}, errors.New("creating server: hostname too short")
  }

  // Return server
  return Server{host, port, secure, pass}, nil
}

// CreateUser creates and returns a user for use in connections.
// Username (uname) and (real)name will default to (nick)name if empty strings.
// Mode must be 0 for none, 4 for +w, 8 for +i, or 12 for +iw.
func CreateUser(nick string, uname string, real string, mode byte) (User,
    error) {
  // Error if nickname is empty
  if len(nick) < 1 {
    return User{}, errors.New("creating user: nickname too short")
  }

  // Error if bitmaskk is not 0, 4, 8, or 12
  if mode != 0 && mode != 4 && mode != 8 && mode != 12 {
    return User{}, errors.New("creating user: mode bitmask invalid")
  }

  // Return user
  return User{nick, uname, real, mode}, nil
}

// EstablishConnection establishes a connection to the specified IRC server
// using the specified user information. It sends initial messages as required
// by the IRC protocol.
func EstablishConnection(Server Server, User User) (Client, error) {
  // Attempt connection establishment
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Server.Host, Server.Port))

  // Return error if one was encountered
  if err != nil {
    return Client{}, err;
  }

  // Create client with connection details
  Client := Client{true, Server, User, conn}

  // Start read loop to get data from server
  go readLoop(Client)

  // Send initial messages to server
  sendRawf(Client, "NICK %s", User.Nick)
  sendRawf(Client, "USER %s %d * :%s", User.User, User.Mode, User.Real)

  // Return client
  return Client, nil;
}


// readLoop is a goroutine used to read data from the server a client is
// connected to.
func readLoop(Client Client) {
  // Create reader to get data from server
  reader := bufio.NewReader(Client.Conn)

  // Loop through reads until client is stopped
  for Client.Active {
    msg, err := reader.ReadString('\n')

    // Output error if encountered
    if err != nil {
      fmt.Printf("[ERROR] %s\n", err);
    }

    // Output data
    fmt.Printf(msg)
  }
}

// sendRawf formats a string to create a raw IRC message that is sent to the
// client's server.
func sendRawf(Client Client, format string, a ...interface{}) {
  fmt.Fprintf(Client.Conn, format + "\r\n", a...);
}

// sendRaw sends a raw IRC message to the client's server.
func sendRaw(Client Client, msg string) {
  fmt.Fprint(Client.Conn, msg + "\r\n");
}
