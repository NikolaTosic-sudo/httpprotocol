package main

import (
	"httpprotocol/internal/request"
	"httpprotocol/internal/response"
	"httpprotocol/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func respond400() []byte {
	return []byte(`
	<html>
	  <head>
		<title>400 Bad Request</title>
	  </head>
	  <body>
	    <h1>Bad Request</h1>
	    <p>Your request honestly kinda sucked.</p>
	  </body>
	</html>
	`)
}

func respond500() []byte {
	return []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`)
}

func respond200() []byte {
	return []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
	`)
}

func main() {
	httpServer, err := server.Serve(port, func(w *response.Writer, req *request.Request) {

		body := respond200()
		statusLine := response.OK

		if req.RequestLine.RequestTarget == "/yourproblem" {
			body = respond400()
			statusLine = response.BadRequest
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			body = respond500()
			statusLine = response.InternalError
		}

		defaultHeaders := response.GetDefaultHeaders(len(body))
		defaultHeaders.Replace("Content-Type", "text/html")
		err := w.WriteStatusLine(statusLine)
		if err != nil {
			log.Fatalf("Error writing status line: %v", err)
		}

		err = w.WriteHeaders(defaultHeaders)
		if err != nil {
			log.Fatalf("Error writing headers: %v", err)
		}

		_, err = w.WriteBody([]byte(body))
		if err != nil {
			log.Fatalf("Error writing body: %v", err)
		}
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
