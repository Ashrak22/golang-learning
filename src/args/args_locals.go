package args

import (
	"os"
	"fmt"
)

type arg struct {
    function ArgFunc
    paramCount int
}

func (a *Argument) stripArgs(arg string) string {
	return arg[len(a.separator):]
}

func (a *Argument) isArg(arg string) bool {
	return arg[0:len(a.separator)] == a.separator
}

func (a *Argument) countMaxArgsAndParams() int {
	count := 0
	for _,item := range a.argsMap {
		count += item.paramCount + 1
	}
	return count
}

func (a *Argument) splitArgs(args []string) map[string][]string {
	res := make(map[string][]string)
	name := ""
	var arr []string
	for _, item := range args {
		if !a.isArg(item) {
			arr = append(arr, item)
		} else {
			if name != "" {
				res[name] = arr
				arr = make([]string, 0)
			}
			name = a.stripArgs(item)
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
