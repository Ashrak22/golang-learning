package args

type arg struct {
	function   ArgFunc
	paramCount int
	separator  string
}

//ArgFunc is type for callbacks on recognized params
type ArgFunc func(...string) error

//Argument class that holds all arguments and needed info
type Argument struct {
	argsMap   map[string]arg
	sepString []byte
}
