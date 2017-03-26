package main

import "net"

const myFacility uint16 = 0x1003

var myErrors = map[uint16]string{
	0x0001: "Couldn't parse/lookup IP: ",
	0x0002: "Port argument is not a valid number",
	0x0003: "Cannot use a Reserved Portnumber",
	0x0004: "Couldn't reach server: ",
	0x0005: "Write failed: ",
	0x0006: "Connection not allowed.",
	0x0007: "Read failed: ",
	0x0008: "Unknown command",
	0x0009: "Marshalling error: ",
	0x0010: "exit",
	0x0011: "Unknown command: %s",
}

var ipaddr []net.IP
var port uint16
var compress = false

var commands = map[string]int32{}
