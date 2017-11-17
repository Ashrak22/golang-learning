package communicator

/********************************************************************************
*
*	TODO:	1) Make writing async (channel many producers one consument)
*				so that we don't have spend as much a time with a locked mutex
*
*
*********************************************************************************/

import (
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
)

//ServerCommunicator wraps network functionality for Server side of the communication, reading, writing, buffered writing
type ServerCommunicator struct {
	conn     *net.TCPConn
	compress bool
	mutex    *sync.Mutex
}

//NewServerCommunicator is the factory function to create new Server Communicator, you have to supply TCP Connection conn
//and bool that specifies whether connection
func NewServerCommunicator(conn *net.TCPConn, compress bool) *ServerCommunicator {
	res := new(ServerCommunicator)
	res.conn = conn
	res.compress = compress
	res.mutex = &sync.Mutex{}
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

func (sc *ServerCommunicator) Lock() {
	sc.mutex.Lock()
}

func (sc *ServerCommunicator) Unlock() {
	sc.mutex.Unlock()
}
