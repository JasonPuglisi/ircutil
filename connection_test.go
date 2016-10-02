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
		{"irc.example.com", 0, false, "", ""},
		{"localhost", 1, true, "Password1!", ""},
		{"10.0.0.1", 65534, false, "", ""},
		{"::1", 65535, false, "", ""},
		{"", 0, false, "", "creating server: hostname too short"},
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

	u := cases[0]
	user, _ := CreateUser(u.nick, u.user, u.real)
	if user.nick != u.nick || user.user != u.user || user.real != u.real {
		t.Errorf("CreateUser(\"%s\", \"%s\", \"%s\") produced overriden values "+
			"{\"%s\", \"%s\", \"%s\"}, no overriden values expected", u.nick, u.user,
			u.real, user.nick, user.user, user.real)
	}

	user, _ = CreateUser(u.nick, "", u.real)
	if user.user != u.nick {
		t.Errorf("CreateUser(\"%s\", \"\", \"%s\") produced an incorrect "+
			"overriden username \"%s\", expected \"%s\"", u.nick, u.real, user.user,
			u.nick)
	}

	user, _ = CreateUser(u.nick, u.user, "")
	if user.real != u.nick {
		t.Errorf("CreateUser(\"%s\", \"%s\", \"\") produced an incorrect "+
			"overriden realname \"%s\", expected \"%s\"", u.nick, u.user, user.real,
			u.nick)
	}
}
