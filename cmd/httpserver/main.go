package main

import (
	"httpprotocol/internal/request"
	"httpprotocol/internal/response"
	"httpprotocol/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	httpServer, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {

		if req.RequestLine.RequestTarget == "/yourproblem" {
			e := server.HandlerError{
				StatusCode: response.BadRequest,
				Message:    []byte("Your problem is not my problem\n"),
			}

			return &e
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			e := server.HandlerError{
				StatusCode: response.InternalError,
				Message:    []byte("Woopsie, my bad\n"),
			}

			return &e
		}

		w.Write([]byte("All good, frfr\n"))

		return nil
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer httpServer.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
