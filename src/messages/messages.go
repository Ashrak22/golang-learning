package messages

import (
	"bettererror"
	"fmt"
	"functions"
	"net"

	"github.com/golang/protobuf/proto"
)

func init() {
	bettererror.RegisterFacility(myFacility, "messages")
}

func WriteMessage(conn *net.TCPConn, msg interface{}) error {
	marshalled, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0003, fmt.Sprintf(myErrors[0x0003], err.Error()))
	}
	_, err = conn.Write(marshalled)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0002, fmt.Sprintf(myErrors[0x0002], err.Error()))
	}
	return nil
}

func ReadMessage(conn *net.TCPConn, msg interface{}, buffer []byte) error {
	functions.Memset(buffer, 0)
	length, err := conn.Read(buffer)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0001, fmt.Sprintf(myErrors[0x0001], err.Error()))
	}
	err = proto.Unmarshal(buffer[:length], msg.(proto.Message))
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0004, fmt.Sprintf(myErrors[0x0004], err.Error()))
	}
	return nil
}
