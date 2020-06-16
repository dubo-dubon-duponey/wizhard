package utils

import (
  "fmt"
  "io"
	"net"
	"time"
)

const maxBufferSize = 1024
const timeout = time.Duration(10000000)

// Uber dirty synchronous UDP client - FIXME
func UDPClient( /*ctx context.Context,*/ address string, reader io.Reader) (string, error) {
  fmt.Println("Opening com with", address)
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return "", err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	for {
		_, err := io.Copy(conn, reader)

		if err != nil {
			return "", err
		}

		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return "", err
		}
		return string(buffer[0:n]), nil
	}

	//	return "", nil
}
