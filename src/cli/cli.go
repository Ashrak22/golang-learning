package main

import (
	"args"
	"bettererror"
	"fmt"
	"messages"
	"net"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
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

func memset(a []byte, v byte) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func writeMessage(conn *net.TCPConn, msg interface{}) error {
	marshalled, err := proto.Marshal(msg.(proto.Message))
	_, err = conn.Write(marshalled)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0005, myErrors[0x0005]+err.Error())
	}
	return nil
}

func readMessage(conn *net.TCPConn, msg interface{}, buffer []byte) error {
	memset(buffer, 0)
	length, err := conn.Read(buffer)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007]+err.Error())
	}
	err = proto.Unmarshal(buffer[:length], msg.(proto.Message))
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0009, myErrors[0x0009]+err.Error())
	}
	return nil
}

func getCommand() (*messages.Command, error) {
	var b = make([]byte, 8*1024)
	fmt.Print("> ")
	memset(b, 0)
	os.Stdin.Read(b)
	var command = string(b)

	trimmed := strings.Trim(command, "\r\n\t "+string(0)) + string(0)
	if strings.HasPrefix(trimmed, "exit") {
		return nil, bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010])
	}

	var comm = &messages.Command{Magic: 0xABCD}
	for key, value := range messages.Commands {
		if strings.HasPrefix(command, key) {
			comm.Command = value
			comm.Argstring = string(command[len(key)+1:])
		}
	}
	if comm.Command == 0 {
		return nil, bettererror.NewBetterError(myFacility, 0x0011, fmt.Sprintf(myErrors[0x0011], trimmed))
	}
	return comm, nil
}

func runLoop() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ipaddr[0], Port: int(port)})
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0004, myErrors[0x0004]+err.Error())
	}
	var buffer = make([]byte, 100*1024)
	defer conn.Close()

	var initMessage = &messages.Init{Version: 1, Magic: 0xABCD, App: "cli"}
	writeMessage(conn, initMessage)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0005, myErrors[0x0005]+err.Error())
	}

	var initResponse = new(messages.InitResponse)
	err = readMessage(conn, initResponse, buffer)
	if !initResponse.Allowed {
		return bettererror.NewBetterError(myFacility, 0x0006, myErrors[0x0006])
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

		writeMessage(conn, command)
		if err != nil {
			return err
		}

		var commandResponse = new(messages.CommandResult)
		readMessage(conn, commandResponse, buffer)
		if err != nil {
			return err
		}
		if commandResponse.CommandResult != 0 {
			fmt.Println(commandResponse.DisplayText)
		}

	}
	return nil
}
