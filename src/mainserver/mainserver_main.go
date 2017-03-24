package main

import (
	"args"
	"bettererror"
	"fmt"
	"messages"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

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
		ret = exec.Command("screen", "-S", "cli", "-p", "0", "-X", "stuff", "cli /port "+strconv.Itoa(int(port))+" \n")
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
		defer conn.Close()
		go runConnection(conn)
	}
	return nil
}

func runConnection(conn *net.TCPConn) {
	var buffer = make([]byte, 100*1024)
	var initmsg = new(messages.Init)

	if err := messages.ReadMessage(conn, initmsg, buffer); err != nil {
		fmt.Println(err.Error())
		return
	}
	if initmsg.Magic != 0xABCD {
		err := bettererror.NewBetterError(myFacility, 0x0009, myErrors[0x0009])
		fmt.Println(err.Error())
		return
	}
	if initmsg.App == "cli" {
		fmt.Println("App cli has connected from ", conn.RemoteAddr().String())
		handleCli(conn, initmsg.Port)
	}
}

func handleCli(conn *net.TCPConn, commPort int32) {
	var initResponse = &messages.InitResponse{Magic: 0xABCD, Allowed: true}

	if err := messages.WriteMessage(conn, initResponse); err != nil {
		fmt.Println(err.Error())
		return
	}

	buffer := make([]byte, 100*1024)
	for true {
		var comm = new(messages.Command)

		if err := messages.ReadMessage(conn, comm, buffer); err != nil {
			fmt.Println(err.Error())
			return
		}

		var resp *messages.CommandResult
		if comm.Magic != 0xABCD {
			err := bettererror.NewBetterError(myFacility, 0x0009, fmt.Sprintf("%s: 0x%4X", myErrors[0x0009], comm.Magic))
			resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: int32(err.Code()), DisplayText: err.Error()}
		} else {
			fmt.Printf("Received command 0x%.8X with args '%s'\r\n", comm.Command, comm.Argstring)
			resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: 0}
		}

		if err := messages.WriteMessage(conn, resp); err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010]+err.Error())
			fmt.Println(err.Error())
			return
		}
	}
}
