package args

import "os"

type arg struct {
    function ArgFunc
    paramCount int
}

func stripArgs(arg string) string {
	return arg[2:len(arg)]
}

func isArg(arg string) bool {
	return arg[0:2] == "--"
}

func (a *Argument) countMaxArgsAndParams() int {
	count := 0
	for _,item := range a.argsMap {
		count += item.paramCount + 1
	}
	return count
}

func splitArgs(args []string) map[string][]string {
	res := make(map[string][]string)
	name := ""
	var arr []string
	for _, item := range args {
		if !isArg(item) {
			arr = append(arr, item)
		} else {
			if name != "" {
				res[name] = arr
				arr = make([]string, 0)
			}
			name = stripArgs(item)
		}
	}
	res[name] = arr
	return res
}

func (a *Argument) raiseError() {
	item, existed := a.argsMap["help"]
	if existed {
		item.function()
	}
	os.Exit(1)
}
