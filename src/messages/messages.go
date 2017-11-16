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

//NewClientStreamCommunicator is factory function for ClientStreamCommunicator
func NewClientStreamCommunicator(port int, host []net.IP, errFunc MsgError, unmarshFunc MsgUnmarshaller, done chan bool) (*ClientStreamCommunicator, error) {
	res := new(ClientStreamCommunicator)
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: host[0], Port: port})
	if err != nil {
		return nil, bettererror.NewBetterError(myFacility, 0x0004, myErrors[0x0004]+err.Error())
	}
	//defer conn.Close()
	res.conn = conn
	res.done = done
	res.errFunc = errFunc
	res.unmarFunc = unmarshFunc
	return res, nil
}

//ReadMessageStream makes async message reading possible
func (c *ClientStreamCommunicator) ReadMessageStream() {
	var channel = make(chan readInfo, 20)
	go func() {
		for {
			var info readInfo
			info.data, info.err = dataReadStream(c.conn)
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
				c.errFunc(info.err)
				continue
			}
			err := c.unmarFunc(info.data)
			if err != nil {
				c.errFunc(err)
			}

		}
		c.done <- true
	}()
}
