package messages

import (
	"net"
)

//MsgUnmarshaller ist the callback for correct Unmarshalling the messages received
type MsgUnmarshaller func(buffer []byte) error

// MsgError is a function to be called onError
type MsgError func(err error)

type readInfo struct {
	data []byte
	err  error
}

//ClientStreamCommunicator is struct that will get methods for stream Reading
type ClientStreamCommunicator struct {
	conn      *net.TCPConn
	done      chan bool
	end       chan bool
	unmarFunc MsgUnmarshaller
	errFunc   MsgError
}
