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

	req, err := request.RequestFromReader(conn)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Request line:")
	fmt.Printf("- Method: %v\n", req.RequestLine.Method)
	fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
	fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)

	fmt.Println("Headers:")
	for k, v := range req.Headers {
		fmt.Printf("- %s: %s\n", k, v)
	}

	fmt.Println("connection has been closed")
}
