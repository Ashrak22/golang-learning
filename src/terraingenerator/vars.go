package main

import "net"

/*Error handling variables and consts*/
const myFacility uint16 = 0x1006

var myErrors = map[uint16]string{
	0x0001: "Port argument is not a valid number",
	0x0002: "Cannot use a Reserved Portnumber",
	0x0003: "Cannot execute screen",
	0x0004: "Cannot execute cli command",
	0x0005: "Listening Port can't be 0",
	0x0006: "Couldn't start listening on Port. ",
	0x0016: "Not implemented",
}

var ipaddr []net.IP
var port uint16
var compress = false
