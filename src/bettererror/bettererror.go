//Package bettererror handlers errors better by adding an error code
package bettererror

//BetterError is the struct that gets returned
type BetterError struct {
	msg  string
	code int
}

//NewBetterError is factory to create new instance of the Better Error struct
func NewBetterError(code int, msg string) *BetterError {
	var err *BetterError
	err = new(BetterError)
	err.msg = msg
	err.code = code
	return err
}

func (e *BetterError) Error() string {
	return e.msg
}

func (e *BetterError) Code() int {
	return e.code
}
