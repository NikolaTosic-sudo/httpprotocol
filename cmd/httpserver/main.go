package main

import (
	"crypto/sha256"
	"fmt"
	"httpprotocol/internal/headers"
	"httpprotocol/internal/request"
	"httpprotocol/internal/response"
	"httpprotocol/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func toStr(bytes []byte) string {
	out := ""

	for _, b := range bytes {
		out += fmt.Sprintf("%02x", b)
	}

	return out
}

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

		defaultHeaders := response.GetDefaultHeaders(len(body))

		if req.RequestLine.RequestTarget == "/yourproblem" {
			body = respond400()
			statusLine = response.BadRequest
			strLen := strconv.Itoa(len(body))
			defaultHeaders.Replace("Content-Length", strLen)
		}

		if req.RequestLine.RequestTarget == "/myproblem" {
			body = respond500()
			statusLine = response.InternalError
			strLen := strconv.Itoa(len(body))
			defaultHeaders.Replace("Content-Length", strLen)
		}

		if req.RequestLine.RequestTarget == "/video" {
			f, _ := os.ReadFile("assets/vim.mp4")
			defaultHeaders.Replace("Content-Type", "video/mp4")
			defaultHeaders.Replace("Content-Length", fmt.Sprintf("%d", len(f)))

			w.WriteStatusLine(response.OK)
			w.WriteHeaders(defaultHeaders)

			w.WriteBody(f)

			return
		}

		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := req.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])

			if err != nil {
				body = respond500()
				statusLine = response.InternalError
				strLen := strconv.Itoa(len(body))
				defaultHeaders.Replace("Content-Length", strLen)
			} else {
				w.WriteStatusLine(response.OK)

				defaultHeaders.Delete("Content-Length")
				defaultHeaders.Set("Transfer-encoding", "chunked")
				defaultHeaders.Replace("Content-Type", "text/plain")

				defaultHeaders.Set("Trailer", "X-Content-SHA256")
				defaultHeaders.Set("Trailer", "X-Content-Length")

				w.WriteHeaders(defaultHeaders)

				fullBody := []byte{}

				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						fmt.Println(err, "err")
						break
					}

					_, err = w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))

					if err != nil {
						fmt.Println(err, "hexa err")
					}

					_, err = w.WriteBody(data[:n])

					if err != nil {
						fmt.Println(err, "body err")
					}

					fullBody = append(fullBody, data[:n]...)

					_, err = w.WriteBody([]byte("\r\n"))

					if err != nil {
						fmt.Println(err, "end err")
					}
				}

				w.WriteBody([]byte("0\r\n"))

				tailers := headers.NewHeaders()

				out := sha256.Sum256(fullBody)

				tailers.Set("X-Content-SHA256", toStr(out[:]))
				tailers.Set("X-Content-Length", strconv.Itoa(len(fullBody)))

				fmt.Println(tailers, "tailers")

				w.WriteHeaders(tailers)

				return
			}
		}

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
