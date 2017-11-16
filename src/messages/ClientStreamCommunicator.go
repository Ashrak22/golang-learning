package messages

import (
	"bettererror"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

//ClientStreamCommunicator is struct that will get methods for stream Reading
type ClientStreamCommunicator struct {
	conn      *net.TCPConn
	done      chan bool
	end       chan bool
	unmarFunc MsgUnmarshaller
	errFunc   MsgError
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

//StartRead starts async reading of messages, which are the passed to unmarshaller after reading
func (c *ClientStreamCommunicator) StartRead() {
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

func (c *ClientStreamCommunicator) Write(msg proto.Message, compress bool) error {
	return WriteMessage(c.conn, msg, compress)
}
