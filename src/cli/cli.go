package main

import (
	"args"
	"bettererror"
	"fmt"
	"functions"
	"messages"
	"os"
	"strings"
	"time"
)

func init() {
	bettererror.RegisterFacility(myFacility, "cliapp")
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
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("IP: %s\r\n", ipaddr[0].String())
	err = runLoop()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getCommand() (*messages.Command, error) {
	var b = make([]byte, 8*1024)
	for {
		time.Sleep(100 * time.Millisecond)
		if internalState == stateIdle {
			break
		}
	}
	fmt.Print("> ")
	functions.Memset(b, 0)
	os.Stdin.Read(b)
	var command = string(b)

	trimmed := strings.Trim(command, "\r\n\t "+string(0))
	if strings.HasPrefix(trimmed, "exit") {
		return nil, bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010])
	}

	var comm = &messages.Command{Magic: 0xABCD}
	for key, value := range commands {
		if strings.HasPrefix(trimmed, key) {
			comm.Command = value
			if len(trimmed) != len(key) {
				comm.Argstring = string(trimmed[len(key)+1:])
			}
		}
	}
	if comm.Command == 0 {
		return nil, bettererror.NewBetterError(myFacility, 0x0011, fmt.Sprintf(myErrors[0x0011], trimmed))
	}
	return comm, nil
}

func runLoop() error {
	done := make(chan bool)
	comm, err := messages.NewClientStreamCommunicator(int(port), ipaddr, errFunc, unmarshaller, done)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0004, myErrors[0x0004]+err.Error())
	}
	defer comm.Close()

	var initMessage = &messages.Init{Version: 1, Magic: 0xABCD, App: "cli", Compress: compress, Port: 40000}
	err = comm.Write(initMessage, false)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0005, myErrors[0x0005]+err.Error())
	}

	for true {
		command, err := getCommand()
		if err != nil {
			if err.(*bettererror.BetterError).Code() == 0x10030010 {
				break
			}
			fmt.Println(err.Error())
			continue
		}

		err = comm.Write(command, compress)
		if err != nil {
			return err
		}
		internalState = stateWaitingForResponse
	}
	return nil
}
