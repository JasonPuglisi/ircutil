// Package ircutil provides utility functions for handling IRC connections and
// operations.
package ircutil

import (
  "bufio"
  "errors"
  "fmt"
  "net"
)

// Struct server stores provided server information.
type server struct {
  host   string
  port   uint16
  secure bool
  pass   string
}

// Struct user stores provided user information, and is not always current.
type user struct {
  nick string
  user string
  real string
  mode byte
}

// CreateServer creates and returns a server struct for use in connections.
// (Host)name is an IP address or FQDN. IPv6 addresses must be in [brackets].
// Secure is true if the specified port uses TLS/SSL.
// (Pass)word must be an empty string if there is no server password.
func CreateServer(host string, port uint16, secure bool, pass string) (server,
    error) {
  if len(host) < 1 {
    return server{}, errors.New("creating server: no hostname")
  }

  return server{host, port, secure, pass}, nil
}

// CreateUser creates and returns a user struct for use in connections.
// Username (uname) and (real)name will default to (nick)name if empty strings.
// Mode must be 0 for none, 4 for +w, 8 for +i, or 12 for +iw.
func CreateUser(nick string, uname string, real string, mode byte) (user,
    error) {
  if len(nick) < 1 {
    return user{}, errors.New("creating user: invalid nickname")
  }

  if len(uname) < 1 {
    return user{}, errors.New("creating user: invalid username")
  }

  if len(real) < 1 {
    return user{}, errors.New("creating user: invalid realname")
  }

  if mode != 0 && mode != 4 && mode != 8 && mode != 12 {
    return user{}, errors.New("creating user: invalid mode bitmask")
  }

  return user{nick, uname, real, mode}, nil
}

// EstablishConnection establishes a connection to the specified IRC server
// using the specified user information.
func EstablishConnection(server server, user user) {
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.host, server.port))

  if err != nil {
    fmt.Printf("Error connecting to server: %s\n", err);
    return
  }

  fmt.Printf("Connected to server %s (%s)\n", server.host, conn.RemoteAddr())

  go readLoop(conn)

  sendRawf(conn, "NICK %s", user.nick)
  sendRawf(conn, "USER %s %d * :%s", user.user, user.mode, user.real)

  for {}
}

func Err() {

}

func readLoop(conn net.Conn) {
  reader := bufio.NewReader(conn)

  fmt.Println("Started read loop")

  for {
    msg, err := reader.ReadString('\n')

    if err != nil {
      fmt.Printf("%s\n", err);
    }

    fmt.Printf(msg)
  }
}

func sendRaw(conn net.Conn, msg string) {
  fmt.Fprint(conn, msg + "\r\n");
}

func sendRawf(conn net.Conn, format string, a ...interface{}) {
  fmt.Fprintf(conn, format + "\r\n", a...);
}
