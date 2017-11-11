package messages

var myFacility uint16 = 0x0002

var myErrors = map[uint16]string{
	0x0001: "Error while reading: %s",
	0x0002: "Error while writing: %s",
	0x0003: "Error while marshalling: %s",
	0x0004: "Error while unmarshaling: %s",
	0x0005: "Error while compressing: %s",
	0x0006: "Error while uncompressing: %s",
	0x0007: "exit",
	0x0008: "Buffer is too small",
}

const bufferSize int = 100 * 1024
