package args

type ArgFunc func(...string)

type Argument struct {
	argsMap map[string]arg
	separator string
}

func NewArg() *Argument {
	argg := new(Argument)
	argg.argsMap = make(map[string]arg)
	argg.separator = "--"
	return argg
}

func (a *Argument) SetSeparator(sep string) {
	a.separator = sep
}


func (a *Argument) RegisterArg(name string, argument ArgFunc, count int) bool {
	_,existed := a.argsMap[name]
	if !existed {
		a.argsMap[name] = arg{argument, count}
		return true
	}
	return false
}

func GetVersion() string{
	return "1.0.0"
}

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
