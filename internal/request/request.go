package request

import (
	"bytes"
	"fmt"
	"io"
)

var (
	ErrorMalformedRequestLine  = fmt.Errorf("malformed request line")
	ErrorUnsuportedHTTPVersion = fmt.Errorf("unsuported http version")
	ErrorRequestInErrorState   = fmt.Errorf("error in request state")

	SEPARATOR = []byte("\r\n")
)

type parserState string

const (
	StateInit  parserState = "initialized"
	StateError parserState = "error"
	StateDone  parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}
type Request struct {
	RequestLine RequestLine
	state       parserState
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, ErrorRequestInErrorState
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.state = StateDone

		case StateDone:
			break outer
		}
	}

	return read, nil
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split((parts[2]), []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen+n])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
