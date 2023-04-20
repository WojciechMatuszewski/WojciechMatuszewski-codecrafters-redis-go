package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	HOST = "0.0.0.0"
	PORT = "6379"
)

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		_, err = conn.Read(make([]byte, 1024))
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}

			fmt.Println("Error reading from connection: ", err.Error())
			os.Exit(1)
		}

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}

		conn.Close()
	}
}
