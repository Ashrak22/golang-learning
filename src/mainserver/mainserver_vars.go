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
}

var port uint16
var cli bool
