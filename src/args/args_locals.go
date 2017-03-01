//Package args provides support for arghument handling
package args

import (
	"os"
	"fmt"
)

type arg struct {
    function ArgFunc
    paramCount int
	separator string
}

func (a *Argument) parseArg(arg string) (bool, string) {
	for key, value := range a.argsMap {
		if arg[:len(value.separator)] == value.separator && arg[len(value.separator):] == key {
			return true, key
		}
	}
	return false, ""
}

func (a *Argument) splitArgs(args []string) map[string][]string {
	res := make(map[string][]string)
	name := ""
	var arr []string
	for _, item := range args {
		isArg, argName := a.parseArg(item)
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
	return res
}

func (a *Argument) raiseError() {
	fmt.Println("Invalid arguments")
	item, existed := a.argsMap["help"]
	if existed {
		item.function()
	}
	os.Exit(1)
}
