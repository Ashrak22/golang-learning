package bettererror

var facilities map[uint16]string

const myFacility uint16 = 0x0000

func init() {
	facilities = make(map[uint16]string)
	facilities[myFacility] = "bettererror"
}

//RegisterFacility saves facility code for resolution in CheckError
func RegisterFacility(facility uint16, name string) error {
	_, existed := facilities[facility]
	if existed {
		return NewBetterError(myFacility, 0x0001, "Facility already exists")
	}
	facilities[facility] = name
	return nil
}

//GetVersion returns the version string
func GetVersion() string {
	return "1.0.1"
}
