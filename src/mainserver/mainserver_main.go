package main

import (
	"args"
	"bettererror"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

/*Error handling variables and consts*/
const myFacility uint16 = 0x1001

var myErrors = map[uint16]string{
	0x0001: "Port argument is not a valid number",
	0x0002: "Cannot use a Reserved Portnumber",
	0x0003: "Cli couldn't be found",
}

var port uint16
var cli bool

func init() {
	bettererror.RegisterFacility(myFacility, "MainServer")
}

func setPort(vs ...string) error {
	i, err := strconv.Atoi(vs[0])
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0001, myErrors[0x0001])
	}
	if i < 1024 {
		return bettererror.NewBetterError(myFacility, 0x0002, myErrors[0x0002])
	}
	port = uint16(i)
	return nil
}

func setCli(vs ...string) error {
	cli = true
	return nil
}

func main() {
	a := args.NewArg()
	err := a.RegisterArg("port", args.ArgFunc(setPort), 1, "/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	err = a.RegisterArg("cli", args.ArgFunc(setCli), 0, "/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	err = a.EvalArgs(os.Args)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if cli {
		_, err := pullUpCli()
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}

}

func pullUpCli() (*exec.Cmd, error) {
	fmt.Println("Pulling up CLI interface")

	var err error
	var ret *exec.Cmd
	//var out []byte

	if runtime.GOOS != "windows" {
		ret = exec.Command("screen", "cli", "/port", strconv.Itoa(int(port)))
		_, err = ret.Output()
		if err.Error() == "exit status 1" {
			err = bettererror.NewBetterError(myFacility, 0x0003, myErrors[0x0003])
		}
	} else {
		fmt.Println("Sorry, i can't pull the CLI up automatically")
	}
	return ret, err
}
