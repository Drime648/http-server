package headers


import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaderParse(t *testing.T){
	//Test: Single Valid Jeader
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFoo:Bar\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers.Get("Host"))
	require.Equal(t, "Bar", headers.Get("Foo"))
	require.Equal(t, 34, n)
	require.True(t, done)


	//Test: Multiple Headers Same Key
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069,localhost:42069", headers.Get("Host"))
	require.True(t, done)


	//Test: Invalid Spacing Header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
	
	//Test: Invalid Spacing Header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)

}
