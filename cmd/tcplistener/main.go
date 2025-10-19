package main

import (
	"fmt"
	"httpprotocol/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connection accepted")

	reqLine, err := request.RequestFromReader(conn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Request line:")
	fmt.Printf("- Method: %v\n", reqLine.RequestLine.Method)
	fmt.Printf("- Target: %v\n", reqLine.RequestLine.RequestTarget)
	fmt.Printf("- Version: %v\n", reqLine.RequestLine.HttpVersion)

	fmt.Println("connection has been closed")
}
