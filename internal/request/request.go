package request

import (
	"bytes"
	"fmt"
	"httpprotocol/internal/headers"
	"io"
	"strings"
)

const crlf = "\r\n"

const bufferSize = 1024

type parserState string

const (
	StateInit                  parserState = "init"
	StateDone                  parserState = "done"
	StateRequestParsingHeaders parserState = "requestStateParsingHeaders"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() *Request {
	return &Request{
		State:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)

	readToIndex := 0

	req := newRequest()

	for !req.done() {
		if readToIndex > 0 {
			readN, err := req.parse(buf[:readToIndex])
			if err != nil {
				return nil, err
			}

			copy(buf, buf[readN:readToIndex])
			readToIndex -= readN

			if readToIndex > 0 || req.done() {
				continue
			}
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}

		readToIndex += n
	}

	return req, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, []byte(crlf))

	if idx == -1 {
		return nil, 0, nil
	}

	reqLineText := b[:idx]
	read := idx + len(crlf)

	reqLine, err := parseLineFromBytes(reqLineText)

	if err != nil {
		return nil, 0, err
	}

	return reqLine, read, nil
}

func parseLineFromBytes(b []byte) (*RequestLine, error) {
	reqLine := bytes.Split(b, []byte(" "))

	if len(reqLine) != 3 {
		return nil, fmt.Errorf("invalid number of arguments in the request line")
	}

	method := string(reqLine[0])
	target := string(reqLine[1])
	ver := bytes.Split(reqLine[2], []byte("/"))

	if method != strings.ToUpper(method) {
		return nil, fmt.Errorf("method has to be all uppercase")
	}

	if len(ver) != 2 {
		return nil, fmt.Errorf("incorrect http version")
	}

	if string(ver[0]) != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version")
	}

	if string(ver[1]) != "1.1" {
		return nil, fmt.Errorf("incorrect http version")
	}

	requestLine := RequestLine{
		HttpVersion:   string(ver[1]),
		RequestTarget: target,
		Method:        method,
	}

	return &requestLine, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.done() {
		return 0, fmt.Errorf("reading from a done state")
	}

	read := 0

	switch r.State {
	case StateInit:
		reqLine, noBytes, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if noBytes == 0 {
			return 0, nil
		}

		r.RequestLine = *reqLine
		read += noBytes

		r.State = StateRequestParsingHeaders

		return read, nil

	case StateRequestParsingHeaders:
		n, done, err := r.Headers.Parse(data[read:])

		if err != nil {
			return 0, err
		}

		read += n

		if done {
			r.State = StateDone
		}

		return read, nil

	case StateDone:
		return 0, fmt.Errorf("reading from a done state")

	default:
		return 0, fmt.Errorf("unknown state")

	}
}

func (r *Request) done() bool {
	return r.State == StateDone
}
