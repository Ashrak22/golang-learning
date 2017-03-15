package main

import (
	"args"
	"bettererror"
	"fmt"
	"os"
)

func version(vs ...string) {
	fmt.Println("Version: 0.0.2-alpha1")
	fmt.Println("Packages:")
	fmt.Println("Args:", args.GetVersion())
	fmt.Println("BetterError:", bettererror.GetVersion())
	//os.Exit(0)
}

func print(vs ...string) {
	fmt.Println(vs[0])
	os.Exit(0)
}

func help(vs ...string) {
	fmt.Println("Usage:")
	fmt.Println("testgo [argument [parameter]]")
	fmt.Println()
	fmt.Println("Possible arguments")
	fmt.Println("-version 				Prints current version, and version of all non-standard packages")
	fmt.Println("-help					This Helptext")
	fmt.Println("-print [text w/o ws]	prints text passed as parameter")
}
