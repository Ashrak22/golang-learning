package args

type ArgError struct {
	msg string
	code int
}

func NewArgError(code int, msg string) ArgError {
	var err ArgError
	err = new(ArgError)
	err.msg = msg
	err.code = code 
	return &ArgError
}

func (e *ArgError) Error() string {
	return e.msg
}

func (e *ArgError) Code() int {
	return e.code
}