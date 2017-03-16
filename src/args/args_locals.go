//Package args provides support for arghument handling
package args

import (
	"bettererror"
	"bytes"
	"strings"
)

const myFacility uint16 = 0x0001

type arg struct {
	function   ArgFunc
	paramCount int
	separator  string
}

func init() {
	bettererror.RegisterFacility(myFacility, "args")
}

func (a *Argument) parseArg(arg string) (bool, string, error) {
	trimmed := strings.TrimLeft(arg, string(a.sepString))
	if arg == trimmed {
		return false, "", nil
	}

	value, existed := a.argsMap[trimmed]
	if !existed {
		return false, "", bettererror.NewBetterError(myFacility, 0x0004, "Found possible not registered argument")
	}

	var buff bytes.Buffer
    buff.WriteString(value.separator)
    buff.WriteString(trimmed)

	if arg == buff.String() {
		return true, trimmed, nil
	} else {
		return false, "", bettererror.NewBetterError(myFacility, 0x0005, "Wrong separator used")
	}
}

func (a *Argument) splitArgs(args []string) (map[string][]string, error) {
	res := make(map[string][]string)
	name := ""
	var arr []string
	for _, item := range args {
		isArg, argName, err := a.parseArg(item)
		if err != nil {
			return nil, err
		}
		if !isArg {
			arr = append(arr, item)
		} else {
			if name != "" {
				res[name] = arr
				arr = make([]string, 0)
			}
			name = argName
		}
	}
	res[name] = arr
	return res, nil
}
