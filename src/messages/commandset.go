package messages

var Commands = map[string]int32{
	//0x0000xxxx - control Commands
	"exit": 0x00000001,

	//0x0001xxxx - config Commands
	"create config":      0x00010001,
	"set default config": 0x00010002,
	"read config":        0x00010003,
	"delete config":      0x00010004,
}
