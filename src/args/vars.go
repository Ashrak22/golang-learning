package args

const myFacility uint16 = 0x0001

var myErrors = map[uint16]string{
	0x0001: "Arg already exists",
	0x0002: "Separator already exists",
	0x0003: "Too many or too few params received",
	0x0004: "Found possible not registered argument",
	0x0005: "Wrong separator used",
}
