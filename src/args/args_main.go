//Package args provides support for arghument handling
package args
/*
 *	TODO:
 *	1) Add multiple possible separators
 *	2) Add support for Flag like args as with ls (ls -l -a / ls -la)
 */
 
 //ArgFunc is type for callbacks on recognized params
type ArgFunc func(...string)

//Argument class that holds all arguments and needed info
type Argument struct {
	argsMap map[string]arg
}

//NewArg is the factory function for creation of basic empty Argument class
//At least one separator has to be added with add Separator
func NewArg() *Argument {
	argg := new(Argument)
	argg.argsMap = make(map[string]arg)
	return argg
}


//RegisterArgs registers callback for recognized argument and number of params it can take
func (a *Argument) RegisterArg(name string, argument ArgFunc, count int, separator string) bool {
	_,existed := a.argsMap[name]
	if !existed {
		a.argsMap[name] = arg{argument, count, separator}
		return true
	}
	return false
}

//GetVersion returns the version string of this package
func GetVersion() string{
	return "1.0.0"
}

//EvalArgs evaluates passed cmd arguments and calls provided callbacks
func (a *Argument) EvalArgs(arg []string) {
	args := a.splitArgs(arg[1:])
	for key, value := range args {
		item, existed := a.argsMap[key]
		if !existed || item.paramCount != len(value) {
			a.raiseError()
		}
		item.function(value...)
	}
}
