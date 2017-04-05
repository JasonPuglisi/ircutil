# ircutil

Utility for connecting to multiple IRC servers and storing connection details.

Developed for use with
[Inami IRC Bot](https://github.com/JasonPuglisi/inami-irc-bot). Check out the
readme and use `go get github.com/jasonpuglisi/inami-irc-bot` to try it out.

## Overview

Creates connections to an unlimited number of IRC networks with support for
server passwords and user NickServ passwords. Runs basic validation on
connection details.

Supports sending various pre-defined IRC message types, as well as raw
messages (see [actions.go](actions.go)).

Implements a command interface to be used with bots or other programs that
can respond to text triggers. Also implements a data storage interface for
bots or other programs to store persistent user data in a variety of scopes.

Has built-in logging functions to distinguish IRC connections/clients (see
[utility.go](utility.go)).
