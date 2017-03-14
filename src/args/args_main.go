//Package args provides support for arghument handling
package args

import (
	BetterError "bettererror"
)

/*
 *	TODO:
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

//RegisterArg registers callback for recognized argument and number of params it can take
func (a *Argument) RegisterArg(name string, argument ArgFunc, count int, separator string) error {
	_, existed := a.argsMap[name]
	if existed {
		return BetterError.NewBetterError(1, "Arg already exists")
	}
	a.argsMap[name] = arg{argument, count, separator}
	return nil
}

//GetVersion returns the version string of this package
func GetVersion() string {
	return "1.0.3"
}

//EvalArgs evaluates passed cmd arguments and calls provided callbacks
func (a *Argument) EvalArgs(arg []string) error {
	if len(arg) == 1 {
		return nil
	}

	args := a.splitArgs(arg[1:])

	for key, value := range args {
		item, existed := a.argsMap[key]
		if !existed {
			return BetterError.NewBetterError(2, "Not existing arguments received")
		} else if item.paramCount != len(value) {
			return BetterError.NewBetterError(3, "Too many params received")
		}
		item.function(value...)
	}
	return nil
}
