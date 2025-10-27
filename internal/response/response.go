package response

import (
	"fmt"
	"httpprotocol/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	OK            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

type writerState string

const (
	StatusLine writerState = "statusLine"
	Headers    writerState = "headers"
	Body       writerState = "body"
)

type Writer struct {
	Writer io.Writer
	State  writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var err error
	switch statusCode {
	case OK:
		_, err = w.Writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case BadRequest:
		_, err = w.Writer.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case InternalError:
		_, err = w.Writer.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		_, err = w.Writer.Write([]byte("HTTP/1.1 500 \r\n"))
	}

	w.State = Headers

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	strLen := strconv.Itoa(contentLen)
	h["Content-Length"] = strLen
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	allHeaders := []byte("")

	for k, v := range headers {
		h := []byte(fmt.Sprintf("%s: %s\r\n", k, v))
		allHeaders = append(allHeaders, h...)
	}

	allHeaders = append(allHeaders, []byte("\r\n")...)

	_, err := w.Writer.Write(allHeaders)

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Writer.Write(p)

	return n, err
}
