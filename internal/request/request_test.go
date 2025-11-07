package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"io"
	// "fmt"
)

type ChunkReader struct {
	data string
	numBytesPerRead int
	pos int
}

func (cr *ChunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos + cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	return n, nil

}

func TestRequestLineParse(t *testing.T) {
	//Test for good GET Request Line
	reader := &ChunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

	//Test for Good GET Request Line with a path
	reader = &ChunkReader{
		data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)


	//Test for invalid number of parts in request line
	reader = &ChunkReader{
		data: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestHeaderParse(t *testing.T) {
	//Test: Standard Headers
	reader := &ChunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	host, _ := r.Headers.Get("Host")
	agent, _ := r.Headers.Get("user-agent")
	accept, _ := r.Headers.Get("accept")
	assert.Equal(t, "localhost:42069", host)
	assert.Equal(t, "curl/7.81.0", agent)
	assert.Equal(t, "*/*", accept)

	//Test: Malformed Header
	reader = &ChunkReader{
		data: "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}


func TestBodyParse(t *testing.T) {
	//Test: Standard Body
	reader := &ChunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
					"Host: localhost:42069\r\n" +
					"Content-Length: 13\r\n" +
					"\r\n" +
					"hello world!\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))


	//Test: Body shorter than reported content length
	reader = &ChunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
					"Host: localhost:42069\r\n" +
					"Content-Length: 20\r\n" +
					"\r\n" +
					"partial content",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.Error(t, err)
}
