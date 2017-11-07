package messages

import proto "github.com/golang/protobuf/proto"

//MsgUnmarshaller ist the callback for correct Unmarshalling the messages received
type MsgUnmarshaller func(buffer []byte, message proto.Message) error
