package main

/*Error handling variables and consts*/
const myFacility uint16 = 0x1001

var myErrors = map[uint16]string{
	0x0001: "Port argument is not a valid number",
	0x0002: "Cannot use a Reserved Portnumber",
	0x0003: "Cannot execute screen",
	0x0004: "Cannot execute cli command",
	0x0005: "Listening Port can't be 0",
	0x0006: "Couldn't start listening on Port. ",
	0x0007: "Error when acceptin connection ",
	0x0008: "Error when waiting for data: ",
	0x0009: "Wrong Magic Number received",
	0x0010: "Error when writing data: ",
	0x0011: "Error while unmarshalling data: ",
	0x0012: "Error while marshalling data: ",
	0x0013: "Can't connect to cli: %s because of: %s",
	0x0014: "App %04X, %s already installed",
}

var port uint16
var cli bool

var climap = map[string]bool{}

var commands = map[string]int16{
	//0x0001xxxx - mainserver Commands
	"create config":      0x0001,
	"set default config": 0x0002,
	"read config":        0x0003,
	"delete config":      0x0004,
	"resend commands":    0x0005,
}

var apps = &installedCommands{installedApps: map[int16]app{}}
