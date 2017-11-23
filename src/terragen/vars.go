package main

/*Error handling variables and consts*/
const myFacility uint16 = 0x1006

var myErrors = map[uint16]string{
	0x0001: "Out of bounds",
}
