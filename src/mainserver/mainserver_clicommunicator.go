package main

import (
	"bettererror"
	"fmt"
	"messages"
	"net"
)

var channels map[string]chan bool

func handleCli(conn *net.TCPConn, commPort int32) {
	var initResponse = &messages.InitResponse{Magic: 0xABCD, Allowed: true}
	err := messages.WriteMessage(conn, initResponse)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buffer := make([]byte, 100*1024)
	for true {
		var comm = new(messages.Command)
		err = messages.ReadMessage(conn, comm, buffer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var resp *messages.CommandResult
		if comm.Magic != 0xABCD {
			err = bettererror.NewBetterError(myFacility, 0x0009, fmt.Sprintf("%s: 0x%4X", myErrors[0x0009], comm.Magic))
			resp = &messages.CommandResult{Magic: 0x0000, CommandResult: int32(err.(*bettererror.BetterError).Code()), DisplayText: err.Error()}
		} else {
			fmt.Printf("Received command 0x%.8X with args '%s'\r\n", comm.Command, comm.Argstring)
			resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: 0}
		}

		err = messages.WriteMessage(conn, resp)
		if err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010]+err.Error())
			fmt.Println(err.Error())
			return
		}
	}
}

func sendCommands(portNumber int32, host string) {

}
