package communicator

import (
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
)

//ServerCommunicator wraps network functionality for Server side of the communication, reading, writing, buffered writing
type ServerCommunicator struct {
	conn         *net.TCPConn
	compress     bool
	mutex        *sync.Mutex
	writeChannel chan proto.Message
	ErrChannel   chan error
	doneChannel  chan bool
}

//NewServerCommunicator is the factory function to create new Server Communicator, you have to supply TCP Connection conn
//and bool that specifies whether connection
func NewServerCommunicator(conn *net.TCPConn) *ServerCommunicator {
	res := new(ServerCommunicator)
	res.conn = conn
	res.compress = false
	res.mutex = &sync.Mutex{}
	res.writeChannel = make(chan proto.Message, 20)
	res.ErrChannel = make(chan error, 4)
	res.doneChannel = make(chan bool)
	go res.writeStream()
	return res
}

func (sc *ServerCommunicator) writeStream() {
	var message proto.Message
	for {
		select {
		case msg := <-sc.writeChannel:
			message = msg
		case <-sc.doneChannel:
			return
		}

		err := writeMessage(sc.conn, message, sc.compress)

		sc.ErrChannel <- err
	}
}

func (sc *ServerCommunicator) Read(pb proto.Message) error {
	return readMessage(sc.conn, pb, sc.compress)
}

func (sc *ServerCommunicator) Write(pb proto.Message) {
	sc.writeChannel <- pb
}

func (sc *ServerCommunicator) Close() {
	sc.doneChannel <- true
	close(sc.writeChannel)
	close(sc.ErrChannel)
	close(sc.doneChannel)
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
