package server

import (
	"fmt"
	"httpprotocol/internal/request"
	"httpprotocol/internal/response"
	"io"
	"log"
	"net"
)

type serverState string

const (
	ServerClosed  serverState = "closed"
	ServerRunning serverState = "running"
)

type Server struct {
	state    serverState
	listener net.Listener
}

type HandlerError struct {
	statusCode int
	message    []byte
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func newServer(l net.Listener) *Server {
	return &Server{
		state:    ServerRunning,
		listener: l,
	}

}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	server, err := net.Listen("tcp", portStr)

	if err != nil {
		return nil, err
	}

	newServer := newServer(server)

	go newServer.listen()

	return newServer, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	s.state = ServerClosed
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	err := response.WriteStatusLine(conn, response.OK)

	if err != nil {
		log.Fatal(err)
	}

	headers := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, headers)

	if err != nil {
		log.Fatal(err)
	}
}
