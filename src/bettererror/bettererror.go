//Package bettererror handlers errors better by adding an error code
package bettererror

//BetterError is the struct that gets returned
type BetterError struct {
	msg      string
	code     uint16
	facility uint16
}

//NewBetterError is factory to create new instance of the Better Error struct
func NewBetterError(facility uint16, code uint16, msg string) *BetterError {
	return &BetterError{msg, code, facility}
}

//Error returns error Message.
func (e *BetterError) Error() string {
	return e.msg
}

//Code returns error code bundled together with facility.
func (e *BetterError) Code() uint32 {
	return ((uint32(e.facility))<<16)&0xFFFF0000 | (uint32(e.code))&0x0000FFFF
}
