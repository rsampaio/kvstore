package protocol

import (
	"errors"
	"strconv"
	"strings"
)

// Parser interface defines what a parser should implement.
type Parser interface {
	// Parse should receive a command string as argument
	// and return an error or nil if successful.
	Parse(string) error
}

// Protocol defines the protocol attributes that will be parsed from text.
type Protocol struct {
	Command       string
	Args          []string
	ReceivesValue bool
}

// Parse receives line string without newline and parse command and arguments.
func (p *Protocol) Parse(line string) error {
	parsed := strings.Split(line, " ")
	if len(parsed) < 1 {
		return errors.New("empty command")
	}

	switch {
	case parsed[0] == "SET":
		if len(parsed[1:]) != 2 {
			return errors.New("set invalid arguments")
		}

		if _, err := strconv.Atoi(parsed[2]); err != nil {
			return errors.New("set invalid size")
		}
		p.ReceivesValue = true

	case parsed[0] == "GET":
		if len(parsed[1:]) != 1 {
			return errors.New("get invalid arguments")
		}
		p.ReceivesValue = false

	case parsed[0] == "DELETE":
		if len(parsed[1:]) != 1 {
			return errors.New("delete invalid arguments")
		}
		p.ReceivesValue = false

	case parsed[0] == "STREAM":
		if len(parsed[1:]) != 0 {
			return errors.New("stream invalid arguments")
		}
		p.ReceivesValue = false

	default:
		return errors.New("invalid command")
	}

	p.Command = parsed[0]
	p.Args = parsed[1:]
	return nil
}
