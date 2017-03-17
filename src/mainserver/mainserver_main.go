package main

import (
	"args"
	"bettererror"
	"fmt"
	"os"
	"strconv"
)

func init() {
	bettererror.RegisterFacility(0x1001, "MainServer")
}

var port uint16

func setPort(vs ...string) {
	i, err := strconv.Atoi(vs[0])
	if err != nil {
		panic(err)
	}
	port = uint16(i)
}

func main() {
	a := args.NewArg()
	err := a.RegisterArg("port", args.ArgFunc(setPort), 1, "/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	err = a.EvalArgs(os.Args)
	panerr := recover()
	if panerr != nil {
		//fmt.Print(panerr.(error).Error())
		os.Exit(1)
	}
}
