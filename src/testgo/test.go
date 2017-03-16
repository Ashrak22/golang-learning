package main

import (
	"args"
	"bettererror"
	"fmt"
	"os"
)

const myFacility uint16 = 0x1000

func init() {
	bettererror.RegisterFacility(myFacility, "testgoApp")
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	a := args.NewArg()
	a.RegisterArg("version", version, 0, "--")
	a.RegisterArg("print", print, 1, "/")
	a.RegisterArg("help", help, 0, "-")
	err := a.EvalArgs(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
