package messages

var myFacility uint16 = 0x0002

var myErrors = map[uint16]string{
	0x0001: "Error while reading: %s",
	0x0002: "Error while writing: %s",
	0x0003: "Error while marshalling: %s",
	0x0004: "Error while unmarshaling: %s",
}

var Commands = map[string]int32{
	//0x0000xxxx - control Commands
	"exit": 0x00000001,

	//0x0001xxxx - config Commands
	"create config":      0x00010001,
	"set default config": 0x00010002,
	"read config":        0x00010003,
	"delete config":      0x00010004,
}
