//Package args provides support for arghument handling
package args

import (
	"bettererror"
	"bytes"
	"strings"
)

var myFacility uint16 = 0x0001
var separators []string

type arg struct {
	function   ArgFunc
	paramCount int
	separator  string
}

func init() {
	bettererror.RegisterFacility(myFacility, "args")
	separators = make([]string, 0)
}

func (a *Argument) parseArg(arg string) (bool, string, error) {
	for key, value := range a.argsMap {
		var buff bytes.Buffer
		buff.WriteString(value.separator)
		buff.WriteString(key)
		if arg == buff.String() {
			return true, key, nil
		} else if strings.HasSuffix(arg, key) {
			return false, "", bettererror.NewBetterError(myFacility, 0x0005, "Wrong separator used")
		}
	}
	for _, sep := range separators {
		if strings.HasPrefix(arg, sep) {
			return false, "", bettererror.NewBetterError(myFacility, 0x0004, "Found possible not registered argument")
		}
	}
	return false, "", nil
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
