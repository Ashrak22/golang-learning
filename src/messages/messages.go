package messages

import (
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
