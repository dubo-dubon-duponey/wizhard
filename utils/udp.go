package utils

import (
	//  "context"
	"fmt"
	"io"
	"net"
	"time"
)

const maxBufferSize = 1024
const timeout = time.Duration(10 * time.Second)

type Result struct {
	Message string
	Error   error
}

func UDPClient( /*ctx context.Context,*/ address string, reader io.Reader) (res string, err error) {
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

	doneChan := make(chan Result)

	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- Result{
				"",
				err,
			}
			return
		}

		fmt.Printf("packet-written: bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)

		deadline := time.Now().Add(timeout)
		//    err = conn.SetWriteDeadline(deadline)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- Result{
				"",
				err,
			}
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- Result{
				"",
				err,
			}
			return
		}

		fmt.Printf("packet-received: bytes=%d from=%s: %s\n",
			nRead, addr.String(), string(buffer[0:nRead]))

		doneChan <- Result{
			string(buffer[0:nRead]),
			nil,
		}
	}()

	var foo Result
	select {
	/*  case <-ctx.Done():
	    fmt.Println("cancelled")
	    err = ctx.Err()*/
	case foo = <-doneChan:
	}

	return foo.Message, foo.Error
}
