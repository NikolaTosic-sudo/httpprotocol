package server

import (
	"fmt"
	"httpprotocol/internal/request"
	"httpprotocol/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    []byte
}

type Handler func(w *response.Writer, req *request.Request)

func newServer(l net.Listener, handler Handler) *Server {
	return &Server{
		listener: l,
		handler:  handler,
	}

}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := fmt.Sprintf(":%d", port)
	server, err := net.Listen("tcp", portStr)

	if err != nil {
		return nil, err
	}

	newServer := newServer(server, handler)

	go newServer.listen()

	return newServer, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		log.Fatal(err)
	}

	writer := &response.Writer{
		Writer: conn,
	}

	s.handler(writer, req)

}
