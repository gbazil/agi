/*
Package agi provides useful utilities for interacting with PBX Asterisk from golang's routines
using AGI protocol.

For a example:

	package main

	import (
		"net"
		"time"
		"github.com/gbazil/agi"
	)

	func handleConnection(c net.Conn) {
		defer c.Close()

		var ret string
		m, _ := agi.ReadMap(c)

		switch m["agi_dnid"] {
		case "1234567890":
			ret = "EXEC TRANSFER \"SIP/abc@somehost.com\""
		case "0123456789":
			ret = "EXEC TRANSFER \"SIP/bca@somehost.com\""
		default:
			ret = "HANGUP"
		}

		agi.WriteLine(c, ret)
		agi.Read(c)
	}

	func main() {
		l, _ := net.Listen("tcp", ":4573")

		defer l.Close()

		for {
			c, _ := l.Accept()

			c.SetDeadline(time.Now().Add(time.Second * 5))
			go handleConnection(c)
		}
	}
*/
package agi

import (
	"net"
	"strings"
)

// Read reads to string s from connection c only once or until timeout (if set for c)
func Read(c net.Conn) (s string, err error) {
	b := make([]byte, 1024)
	var n int

	n, err = c.Read(b)
	if err == nil {
		s = string(b[:n])
	}

	return
}

// ReadLines collects input into string s from connection c until meets empty line or timeout occurred (if it set for c)
func ReadLines(c net.Conn) (s string, err error) {
	b := make([]byte, 1024)
	var n int

	for {
		n, err = c.Read(b)
		if err != nil {
			break
		}

		s += string(b[:n])

		if strings.HasSuffix(s, "\n\n") {
			break
		}
	}

	return
}

// Parse parses text (AGI vars) into map m and return it
func Parse(s string) (m map[string]string) {
	m = make(map[string]string)
	for _, val := range strings.Split(s, "\n") {
		pair := strings.Split(val, ": ")
		if len(pair) == 2 {
			m[pair[0]] = pair[1]
		}
	}

	return
}

// ReadMap read agi input into map m from connection c and return it
func ReadMap(c net.Conn) (m map[string]string, err error) {
	var s string
	s, err = ReadLines(c)
	m = Parse(s)
	return
}

// Write writes to connection c string s
func Write(c net.Conn, s string) (n int, err error) {
	n, err = c.Write([]byte(s))

	return
}

// WriteLine writes to connection c string s with NL character
func WriteLine(c net.Conn, s string) (n int, err error) {
	n, err = c.Write([]byte(s + "\n"))

	return
}
