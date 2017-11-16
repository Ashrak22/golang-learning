package messages

//MsgUnmarshaller ist the callback for correct Unmarshalling the messages received
type MsgUnmarshaller func(buffer []byte) error

// MsgError is a function to be called onError
type MsgError func(err error)

type readInfo struct {
	data []byte
	err  error
}
