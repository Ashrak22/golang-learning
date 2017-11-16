package messages

import (
	"bettererror"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

//WriteMessage write message to TCPConnection and if needed compresses it
func WriteMessage(conn *net.TCPConn, msg proto.Message, compress bool) error {
	if compress {
		return writeMessageCompressed(conn, msg)
	}
	return writeMessageUncompressed(conn, msg)
}

//ReadMessage reads message from TCPConnection and if needed uncompresses it
func ReadMessage(conn *net.TCPConn, msg proto.Message, compress bool) error {
	if compress {
		return readMessageCompressed(conn, msg)
	}
	return readMessageUncompressed(conn, msg)
}

//ReadMessageStream makes async message reading possible
func ReadMessageStream(unmarshaller MsgUnmarshaller, errFunc MsgError, conn *net.TCPConn, done chan bool) {
	var channel = make(chan readInfo, 20)
	go func() {
		for {
			var info readInfo
			info.data, info.err = dataReadStream(conn)
			if info.err.(*bettererror.BetterError).Code() == 0x00020007 {
				fmt.Println("Client disconnected")
				channel <- info
				break
			}
			channel <- info
		}
		close(channel)
	}()
	go func() {
		for {
			info, more := <-channel
			if !more {
				break
			}
			if info.err != nil {
				errFunc(info.err)
				continue
			}
			err := unmarshaller(info.data)
			if err != nil {
				errFunc(err)
			}

		}
		done <- true
	}()
}
