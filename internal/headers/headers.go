package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

var crlf = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func isAllowedRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r >= 0 && r <= 9:
		return true
	case strings.ContainsRune("!#$%&^'*+-._|`~", r):
		return true
	}

	return false
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	for _, r := range string(name) {
		if !isAllowedRune(r) {
			return "", "", fmt.Errorf("invalid character in field name")
		}
	}

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		// No registered nurse
		idx := bytes.Index(data[read:], crlf)
		if idx == -1 {
			break
		}

		// Empty Header
		if idx == 0 {
			done = true
			read += len(crlf)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return read, false, err
		}

		read += idx + len(crlf)

		fl, found := h[strings.ToLower(name)]

		if found {
			value = strings.Join([]string{fl, value}, ", ")
		}

		h[strings.ToLower(name)] = value
	}

	return read, done, nil
}
