package bettererror

var facilities map[uint16]string

func init() {
	facilities = make(map[uint16]string)
}

//RegisterFacility saves facility code for resolution in CheckError
func RegisterFacility(facility uint16, name string) {
	facilities[facility] = name
}

func GetVersion() string {
	return "1.0.0"
}
