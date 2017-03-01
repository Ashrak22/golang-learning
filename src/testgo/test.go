package main

import (
	"os"
	"args"
)

func main(){
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	a := args.NewArg()
	a.SetSeparator("/")
	a.RegisterArg("version", version, 0)
	a.RegisterArg("print", print, 1)
	a.RegisterArg("help", help, 0)
	a.EvalArgs(os.Args)
}
