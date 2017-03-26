package messages

import (
	"bettererror"
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
)

func init() {
	bettererror.RegisterFacility(myFacility, "messages")
}

func WriteMessage(conn *net.TCPConn, msg proto.Message, compress bool) error {
	if compress {
		return writeMessageCompressed(conn, msg)
	}
	return writeMessageUncompressed(conn, msg)
}

func ReadMessage(conn *net.TCPConn, msg proto.Message, compress bool) error {
	if compress {
		return readMessageCompressed(conn, msg)
	}
	return readMessageUncompressed(conn, msg)
}

func writeMessageUncompressed(conn *net.TCPConn, msg proto.Message) error {
	marshalled, err := proto.Marshal(msg)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0003, fmt.Sprintf(myErrors[0x0003], err.Error()))
	}
	_, err = conn.Write(marshalled)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0002, fmt.Sprintf(myErrors[0x0002], err.Error()))
	}
	return nil
}

func readMessageUncompressed(conn *net.TCPConn, msg proto.Message) error {
	var buffer = make([]byte, bufferSize)
	length, err := conn.Read(buffer)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007])
		}
		return bettererror.NewBetterError(myFacility, 0x0001, fmt.Sprintf(myErrors[0x0001], err.Error()))
	}

	err = proto.Unmarshal(buffer[:length], msg)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0004, fmt.Sprintf(myErrors[0x0004], err.Error()))
	}
	return nil
}

func writeMessageCompressed(conn *net.TCPConn, msg proto.Message) error {
	var b bytes.Buffer
	marshalled, err := proto.Marshal(msg)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0003, fmt.Sprintf(myErrors[0x0003], err.Error()))
	}

	//compress message
	compressor, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	_, err = compressor.Write(marshalled)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0005, fmt.Sprintf(myErrors[0x0005], err.Error()))
	}
	compressor.Flush()
	compressor.Close()

	_, err = conn.Write(b.Bytes())
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0002, fmt.Sprintf(myErrors[0x0002], err.Error()))
	}

	return nil
}

func readMessageCompressed(conn *net.TCPConn, msg proto.Message) error {
	var localBuff = make([]byte, 2*bufferSize)
	var buffer = make([]byte, bufferSize)

	length, err := conn.Read(buffer)
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007])
		}
		return bettererror.NewBetterError(myFacility, 0x0001, fmt.Sprintf(myErrors[0x0001], err.Error()))
	}

	//uncompress
	b := bytes.NewReader(buffer[:length])
	uncompressor, _ := gzip.NewReader(b)
	length, err = uncompressor.Read(localBuff)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0006, fmt.Sprintf(myErrors[0x0006], err.Error()))
	}

	err = proto.Unmarshal(localBuff[:length], msg)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0004, fmt.Sprintf(myErrors[0x0004], err.Error()))
	}

	return nil
}
