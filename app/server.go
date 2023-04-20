package main

import (
	"bufio"
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
			fmt.Errorf("failed to accept connection: %w", err)
			continue
		}

		fmt.Println("New connection to the server", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) error {
	for {
		_, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// The EOF here is expected.
			// Redis seem to first send the command, then, in the next message, the EOF
			// We cannot use ioutil.ReadAll here as the initial message does not contain EOF
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("failed to read from connection: %w", err)
		}

		conn.Write([]byte("+PONG\r\n"))
	}

	return nil
}
