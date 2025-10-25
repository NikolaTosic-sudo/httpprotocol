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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case BadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case InternalError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		_, err := w.Write([]byte("HTTP/1.1 500 \r\n"))
		return err
	}
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	strLen := strconv.Itoa(contentLen)
	h["Content-Length"] = strLen
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	allHeaders := []byte("")

	for k, v := range headers {
		h := []byte(fmt.Sprintf("%s: %s\r\n", k, v))
		allHeaders = append(allHeaders, h...)
	}

	allHeaders = append(allHeaders, []byte("\r\n")...)

	_, err := w.Write(allHeaders)

	return err
}
