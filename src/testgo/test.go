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
	a := args.NewArg()
	a.RegisterArg("version", version, 0, "--")
	a.RegisterArg("print", print, 1, "/")
	a.RegisterArg("help", help, 0, "-")
	err := a.EvalArgs(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
	for true {
	}
}
