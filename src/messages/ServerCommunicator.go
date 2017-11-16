package messages

import "net"
import "github.com/golang/protobuf/proto"

type ServerCommunicator struct {
	conn     *net.TCPConn
	compress bool
}

func NewServerCommunicator(conn *net.TCPConn, compress bool) *ServerCommunicator {
	res := new(ServerCommunicator)
	res.conn = conn
	res.compress = compress
	return res
}

func (sc *ServerCommunicator) Read(pb proto.Message) error {
	return readMessage(sc.conn, pb, sc.compress)
}

func (sc *ServerCommunicator) Write(pb proto.Message) error {
	return writeMessage(sc.conn, pb, sc.compress)
}

func (sc *ServerCommunicator) Close() {
	sc.conn.Close()
}

func (sc *ServerCommunicator) GetRemoteAddress() string {
	return sc.conn.RemoteAddr().String()
}

func (sc *ServerCommunicator) GetLocalPort() int {
	return sc.conn.LocalAddr().(*net.TCPAddr).Port
}

func (sc *ServerCommunicator) SetCompress(compress bool) {
	sc.compress = compress
}
