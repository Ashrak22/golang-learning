package main

import (
	"args"
	"bettererror"
	"fmt"
	"messages"
	"net"
	"os"
)

func init() {
	bettererror.RegisterFacility(myFacility, "MainServer")
	apps.registerApp(0x0001, "mainserver", commands)
}

func main() {
	a := args.NewArg()
	err := a.RegisterArg("port", args.ArgFunc(setPort), 1, "/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	err = a.RegisterArg("cli", args.ArgFunc(pullUpCli), 0, "/")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	err = a.EvalArgs(os.Args)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if port == 0 {
		err = bettererror.NewBetterError(myFacility, 0x0005, myErrors[0x0005])
		fmt.Print(err.Error())
		os.Exit(1)
	}
	err = runLoop()
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

func runLoop() error {
	fmt.Println("Starting networking subsystem")
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: int(port)})
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0006, myErrors[0x0006]+err.Error())
	}
	defer listener.Close()
	for true {
		fmt.Println("Waiting for connection")
		conn, err := listener.AcceptTCP()
		if err != nil {
			return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007]+err.Error())
		}
		comm := messages.NewServerCommunicator(conn, false)
		go runConnection(comm)
	}
	return nil
}

func runConnection(comm *messages.ServerCommunicator) {
	var initmsg = new(messages.Init)

	if err := comm.Read(initmsg); err != nil {
		fmt.Println(err.Error())
		return
	}
	if initmsg.Magic != 0xABCD {
		err := bettererror.NewBetterError(myFacility, 0x0009, myErrors[0x0009])
		fmt.Println(err.Error())
		return
	}
	if initmsg.App == "cli" {
		fmt.Println("App cli has connected from ", comm.GetRemoteAddress())
		comm.SetCompress(initmsg.Compress)
		handleCli(comm)
	}
}
