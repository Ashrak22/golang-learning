package main

import (
	"args"
	"bettererror"
	"fmt"
	"os"
)

func init() {
	bettererror.RegisterFacility(myFacility, "Terrain Generator")
}

func main() {
	a := args.NewArg()
	err := a.RegisterArg("host", args.ArgFunc(getHost), 1, "--")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = a.RegisterArg("port", args.ArgFunc(setPort), 1, "--")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = a.RegisterArg("compression", args.ArgFunc(setCompression), 1, "--")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = a.EvalArgs(os.Args)
}
