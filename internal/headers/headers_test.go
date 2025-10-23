package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsingHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:     barbar     \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "barbar", headers["foofoo"])
	assert.Equal(t, "", headers["MissingKey"])
	assert.Equal(t, 50, n)
	assert.True(t, done)

	// Test: Valid header with capital letters
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nFOoFoO:     barbar     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "barbar", headers["foofoo"])
	assert.Equal(t, "", headers["MissingKey"])
	assert.Equal(t, 50, n)
	assert.True(t, done)

	// Test: Valid header with multiple same field names
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost:     barbar     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, barbar", headers["host"])
	assert.Equal(t, "", headers["MissingKey"])
	assert.Equal(t, 48, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in the header
	headers = NewHeaders()
	data = []byte("H®st : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
