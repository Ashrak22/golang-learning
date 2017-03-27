package main

import (
	"bettererror"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

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

func pullUpCli(vs ...string) error {
	fmt.Println("Pulling up CLI interface")

	var err error
	var ret *exec.Cmd
	if runtime.GOOS != "windows" {
		ret = exec.Command("screen", "-dmS", "cli", "bash")
		ret.Start()
		time.Sleep(2 * time.Second)
		_, err = ret.Output()
		if err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0003, myErrors[0x0003])
		}
		ret = exec.Command("screen", "-S", "cli", "-p", "0", "-X", "stuff", "cli --port "+strconv.Itoa(int(port))+" --host localhost --compression true \n")
		_, err = ret.Output()
		if err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0004, myErrors[0x0004])
		}
		fmt.Println("CLI succesfully pulled up, you can acces it by executing 'screen -r cli'")
	} else {
		fmt.Println("Sorry, i can't pull the CLI up automatically on windows")
	}

	return err
}
