package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer file.Close()
		defer close(out)

		currLine := ""
		for {
			data := make([]byte, 8)
			n, err := file.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				currLine += string(data[:i])
				data = data[i+1:]
				out <- currLine
				currLine = ""
			}
			currLine += string(data)
		}
		if len(currLine) != 0 {
			out <- currLine
		}

	}()
	return out

}
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("Connection has been accepted")

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		// fmt.Println("Connection closed")
	}

}
