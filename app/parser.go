package main

import (
	"bytes"
	"strings"
)

type DataType string

const (
	DATATYPE_ARRAY = DataType("*")
)

type Command string

const (
	COMMAND_ECHO    = Command("echo")
	COMMAND_PING    = Command("ping")
	COMMAND_UNKNOWN = Command("unknown")
)

type ParsedInput struct {
	DataType DataType
	Command  Command
	Payload  []byte
}

func Parse(rawBuf []byte) (ParsedInput, error) {
	buf := bytes.Trim(rawBuf, "\x00")
	command := string(buf[0])

	var parsedInput ParsedInput

	switch command {
	case string(DATATYPE_ARRAY):
		parsedInput.DataType = DATATYPE_ARRAY

		lines := strings.Split(string(buf), "\r\n")
		for i, line := range lines {
			if i == 0 {
				continue
			}

			if strings.HasPrefix(line, "$") {
				continue
			}

			if line == "" {
				continue
			}

			if line == string(COMMAND_ECHO) {
				parsedInput.Command = COMMAND_ECHO
				continue
			}
			if line == string(COMMAND_PING) {
				parsedInput.Command = COMMAND_PING
				continue
			}

			parsedInput.Payload = append(parsedInput.Payload, []byte(line)...)
		}
	default:
		parsedInput.Command = COMMAND_UNKNOWN
	}

	return parsedInput, nil
}
