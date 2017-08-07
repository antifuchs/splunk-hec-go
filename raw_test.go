package hec

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testReader struct {
	input <-chan []byte
	last  []byte
}

func newTestReader(ch <-chan []byte) *testReader {
	return &testReader{ch, nil}
}

func (tr *testReader) Read(buf []byte) (int, error) {
	if tr.last != nil {
		defer func() {
			tr.last = nil
		}()
		return copy(buf, tr.last), nil
	}
	select {
	case in, ok := <-tr.input:
		if ok {
			n := copy(buf, in)
			if n < len(in) {
				tr.last = in[n:]
			}
			return n, nil
		}
		return 0, io.EOF
	}
}

func TestRawReader(t *testing.T) {
	ch := make(chan []byte)
	tr := newRawReader(newTestReader(ch))
	go func() {
		ch <- []byte{
			'1', '2', '3', '4', '\n',
			'5',
		}
	}()
	buf := make([]byte, 8)
	read, err := tr.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 5, read)
	assert.Equal(t, "1234\n", string(buf[0:read]))

	go func() {
		ch <- []byte{
			'6', '7', '8', '9', '\n',
		}
	}()
	buf = make([]byte, 7)
	read, err = tr.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, 6, read)
	assert.Equal(t, "56789\n", string(buf[0:read]))

	go func() {
		close(ch)
	}()
	read, err = tr.Read(buf)
	assert.Equal(t, io.EOF, err, "Should return EOF")
	assert.Equal(t, 0, read, "Should not have read any bytes")
}
