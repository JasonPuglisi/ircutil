package ircutil

import "testing"

func TestCreateServer(t *testing.T) {
  cases := []struct {
    host string
    port uint16
    secure bool
    pass string
    err string
  }{
    {"irc.example.com", 0, false, "", ""},
    {"localhost", 1, true, "Password1!", ""},
    {"10.0.0.1", 65534, false, "", ""},
    {"::1", 65535, false, "", ""},
    {"", 0, false, "", "creating server: no hostname"},
  }

  for _, c := range cases {
    _, err := CreateServer(c.host, c.port, c.secure, c.pass);
    if (err == nil) && len(c.err) > 0 {
      t.Errorf("CreateServer(\"%s\", %d, %t, \"%s\") did not produce an " +
        "error, expected \"%s\"", c.host, c.port, c.secure, c.pass, c.err)
    }

    if (err != nil) && len(c.err) < 1 {
      t.Errorf("CreateServer(\"%s\", %d, %t, \"%s\") produced an error, " +
        "expected no error", c.host, c.port, c.secure, c.pass)
    }

    if (err != nil) && (err.Error() != c.err) {
      t.Errorf("CreateServer(\"%s\", %d, %t, \"%s\") produced an incorrect " +
        "error \"%s\", expected \"%s\"", c.host, c.port, c.secure, c.pass, err, c.err)
    }
  }
}
