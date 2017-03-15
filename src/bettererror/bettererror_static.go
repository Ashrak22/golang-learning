package bettererror

import (
	"fmt"
)

var facilities map[uint16]string

func init() {
	facilities = make(map[uint16]string)
}

func RegisterFacility(facility uint16, name string) {
	facilities[facility] = name
}

func CheckError(err *BetterError) {
	if err != nil {
		fmt.Printf("Package: %s \r\nErrorcode: 0x%08x\r\nMessage: %s ", facilities[err.facility], err.Code(), err.msg)
	}
}

func GetVersion() string {
	return "1.0.0"
}
