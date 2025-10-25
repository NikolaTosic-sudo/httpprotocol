package server

import (
	"bytes"
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
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    []byte
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func newServer(l net.Listener, handler Handler) *Server {
	return &Server{
		state:    ServerRunning,
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

	req, err := request.RequestFromReader(conn)

	if err != nil {
		log.Fatal(err)
	}

	buff := bytes.NewBuffer([]byte{})

	handlerError := s.handler(buff, req)

	b := buff.Bytes()

	headers := response.GetDefaultHeaders(len(b))

	if handlerError != nil {
		err = response.WriteStatusLine(conn, response.StatusCode(handlerError.StatusCode))

		if err != nil {
			log.Fatal(err)
		}

		headers := response.GetDefaultHeaders(len(handlerError.Message))

		err = response.WriteHeaders(conn, headers)
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write(handlerError.Message)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err = response.WriteStatusLine(conn, response.OK)

	if err != nil {
		log.Fatal(err)
	}

	err = response.WriteHeaders(conn, headers)

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(b)

	if err != nil {
		log.Fatal(err)
	}

}
