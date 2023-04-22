package main

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	HOST = "0.0.0.0"
	PORT = "6379"
)

func main() {
	address := fmt.Sprintf("%s:%s", HOST, PORT)
	server := NewServer(address)
	// This feels bad
	defer server.Stop()

	server.Start()
}

type Server struct {
	listener net.Listener
	address  string
	quitch   chan struct{}
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		quitch:  make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %w", s.address, err)
	}

	s.listener = ln
	defer s.listener.Close()

	go s.acceptLoop()

	<-s.quitch

	return nil
}

func (s *Server) Stop() {
	s.quitch <- struct{}{}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v", err)
			continue
		}

		fmt.Println("New connection to the server", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) error {
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			// The EOF here is expected.
			// Redis seem to first send the command, then, in the next message, the EOF
			// We cannot use ioutil.ReadAll here as the initial message does not contain EOF
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("failed to read from connection: %w", err)
		}
		parsedInput, err := Parse(buf)
		if err != nil {
			fmt.Println("error", err)
		}

		if parsedInput.Command == COMMAND_PING {
			conn.Write([]byte("+PONG\r\n"))
			continue
		}

		if parsedInput.Command == COMMAND_ECHO {
			response := fmt.Sprintf("$%d\r\n%s\r\n", len(parsedInput.Payload), string(parsedInput.Payload))
			conn.Write([]byte(response))

			continue
		}

		fmt.Println("unknown command for input:", string(buf))
	}

	return nil
}
