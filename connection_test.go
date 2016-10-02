package ircutil

import "testing"

func TestCreateServer(t *testing.T) {
	cases := []struct {
		host   string
		port   uint16
		secure bool
		pass   string
		err    string
	}{
		{"irc.example.com", 1, false, "", ""},
		{"localhost", 2, true, "Password1!", ""},
		{"10.0.0.1", 65534, false, "", ""},
		{"::1", 65535, false, "", ""},
		{"", 1, false, "", "creating server: hostname too short"},
		{"irc.example.com", 0, false, "", "creating server: port is zero"},
	}

	for _, c := range cases {
		_, err := CreateServer(c.host, c.port, c.secure, c.pass)

		if err == nil && len(c.err) > 0 {
			t.Errorf("CreateServer(\"%s\", %d, %t, \"%s\") did not produce an "+
				"error, expected \"%s\"", c.host, c.port, c.secure, c.pass, c.err)
		}

		if err != nil && len(c.err) < 1 {
			t.Errorf("CreateServer(\"%s\", %d, %t, \"%s\") produced an error, "+
				"expected no error", c.host, c.port, c.secure, c.pass)
		}

		if err != nil && err.Error() != c.err {
			t.Errorf("CreateServer(\"%s\", %d, %t, \"%s\") produced an incorrect "+
				"error \"%s\", expected \"%s\"", c.host, c.port, c.secure, c.pass, err,
				c.err)
		}
	}
}

func TestCreateUser(t *testing.T) {
	cases := []struct {
		nick string
		user string
		real string
		err  string
	}{
		{"Inami", "inami", "Mahiru Inami", ""},
		{"I", "i", "I", ""},
		{"Inami", "inami", "Mahiru Inami", ""},
		{"Inami", "inami", "Mahiru Inami", ""},
		{"Inami", "inami", "Mahiru Inami", ""},
		{"", "inami", "Mahiru Inami", "creating user: nickname too short"},
	}

	for _, c := range cases {
		_, err := CreateUser(c.nick, c.user, c.real)

		if err == nil && len(c.err) > 0 {
			t.Errorf("CreateUser(\"%s\", \"%s\", \"%s\" did not produce an error, "+
				"expected \"%s\"", c.nick, c.user, c.real, c.err)
		}

		if err != nil && len(c.err) < 1 {
			t.Errorf("CreateUser(\"%s\", \"%s\", \"%s\") produced an error, "+
				"expected no error", c.nick, c.user, c.real)
		}

		if err != nil && err.Error() != c.err {
			t.Errorf("CreateUser(\"%s\", \"%s\", \"%s\") produced an incorrect "+
				"error \"%s\", expected \"%s\"", c.nick, c.user, c.real, err, c.err)
		}
	}
}
