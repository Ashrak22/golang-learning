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

func runLoop() error {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ipaddr[0], Port: int(port)})
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0004, myErrors[0x0004]+err.Error())
	}
	defer conn.Close()
	var initMessage = &messages.Init{Version: 1, Magic: 0xABCD, App: "cli"}
	data, err := proto.Marshal(initMessage)
	_, err = conn.Write(data)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0005, myErrors[0x0005]+err.Error())
	}
	memset(data, 0)
	_, err = conn.Read(data)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007]+err.Error())
	}
	var initResponse = new(messages.InitResponse)
	err = proto.Unmarshal(data, initResponse)
	if !initResponse.Allowed {
		return bettererror.NewBetterError(myFacility, 0x0006, myErrors[0x0006])
	}
	var b = make([]byte, 1024)
	for true {
		fmt.Print("> ")
		memset(b, 0)
		os.Stdin.Read(b)
		var command = string(b)
		command = strings.TrimRight(command, "\r\n")
		if strings.HasPrefix(command, "exit") {
			break
		}
		var comm = &messages.Command{Magic: 0xABCD}
		for key, value := range messages.Commands {
			if strings.HasPrefix(command, key) {
				comm.Command = value
				comm.Argstring = string(command[len(key)+1:])
			}
		}
		if comm.Command == 0 {
			fmt.Println(bettererror.NewBetterError(myFacility, 0x0008, myErrors[0x0008]).Error())
			continue
		}
		memset(data, 0)
		data, err = proto.Marshal(comm)
		_, err = conn.Write(data)
		if err != nil {
			return bettererror.NewBetterError(myFacility, 0x0005, myErrors[0x0005]+err.Error())
		}
		memset(data, 0)
		_, err = conn.Read(data)
		if err != nil {
			return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007]+err.Error())
		}
		var commandResponse = new(messages.CommandResult)
		err = proto.Unmarshal(data, commandResponse)
		if commandResponse.CommandResult != 0 {
			fmt.Println(commandResponse.DisplayText)
		}
	}
	return nil
}
