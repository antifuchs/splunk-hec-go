package hec

import (
	"io"
)

type rawReader struct {
	in            io.Reader
	lastRemainder []byte
}

func newRawReader(in io.Reader) *rawReader {
	return &rawReader{
		in: in,
	}
}

func (rr *rawReader) Read(buf []byte) (int, error) {
	var start int
	if rr.lastRemainder != nil && len(rr.lastRemainder) > 0 {
		start = copy(buf, rr.lastRemainder)
		if len(rr.lastRemainder) > len(buf) {
			// Pathological case: The last line fragment
			// is longer than the buffer. We can only it
			// and carry on sending the remaining line in
			// the next buffer.
			rr.lastRemainder = rr.lastRemainder[start:]
			return start, nil
		}
		rr.lastRemainder = nil
	}
	// fill our buffer:
	read, err := rr.in.Read(buf[start:])
	if err != nil {
		if err == io.EOF && start > 0 {
			return start, nil
		}
		return start + read, err
	}
	// split off the last non-newline portion:
	for i := start + read - 1; i >= start; i-- {
		if buf[i] == '\n' {
			rr.lastRemainder = buf[i+1 : start+read]
			return i + 1, nil
		}
	}
	return start + read, err
}
