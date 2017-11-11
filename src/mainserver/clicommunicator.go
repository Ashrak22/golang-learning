package main

import (
	"bettererror"
	"fmt"
	"messages"
	"net"
	"strconv"
	"strings"
	"time"
)

func handleCli(conn *net.TCPConn, commPort int, compress bool) {
	var initResponse = &messages.InitResponse{Magic: 0xABCD, Allowed: true}
	var name string
	if strings.Contains(conn.RemoteAddr().(*net.TCPAddr).IP.String(), ":") {
		name = "[" + conn.RemoteAddr().(*net.TCPAddr).IP.String() + "]:" + strconv.Itoa(commPort)
	} else {
		name = conn.RemoteAddr().(*net.TCPAddr).IP.String() + ":" + strconv.Itoa(commPort)
	}
	climap[name] = compress
	defer removeConnection(name)
	go sendCommands()
	if err := messages.WriteMessage(conn, initResponse, compress); err != nil {
		fmt.Println(err.Error())
		return
	}

	for true {
		var comm = new(messages.Command)

		if err := messages.ReadMessage(conn, comm, compress); err != nil {
			if err.(*bettererror.BetterError).Code() == 0x00020007 {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println(err.Error())
			return
		}

		var resp *messages.CommandResult
		if comm.Magic != 0xABCD {
			err := bettererror.NewBetterError(myFacility, 0x0009, fmt.Sprintf("%s: 0x%4X", myErrors[0x0009], comm.Magic))
			resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: int32(err.Code()), DisplayText: err.Error()}
		} else {
			var err error
			fmt.Printf("Received command 0x%.8X with args '%s'\r\n", comm.Command, comm.Argstring)
			switch comm.Command {
			case 0x0005:
				err = resendCommands()
			default:
				err = nil
			}
			if err != nil {
				resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: int32(err.(*bettererror.BetterError).Code()), DisplayText: err.Error()}
			} else {
				resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: 0}
			}
		}

		if err := messages.WriteMessage(conn, resp, compress); err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010]+err.Error())
			fmt.Println(err.Error())
			return
		}
	}
}

func resendCommands() error {
	sendCommands()
	return nil
}

func sendCommands() {
	time.Sleep(100 * time.Millisecond)
	for key, value := range climap {
		conn, err := net.Dial("tcp", key)
		if err != nil {
			fmt.Println(bettererror.NewBetterError(myFacility, 0x0013, fmt.Sprintf(myErrors[0x0013], key, err.Error())).Error())
			continue
		}
		var msg = new(messages.CommandPush)
		msg.Magic = 0xABCD
		msg.Commands = apps.flattenCommands()
		if err = messages.WriteMessage(conn.(*net.TCPConn), msg, value); err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

}

func removeConnection(conn string) {
	delete(climap, conn)
}

func commandDispatch(code int32, vs ...string) error {
	appfacility := int16(code >> 16)
	value, existed := apps.installedApps[appfacility]
	if !existed {
		return bettererror.NewBetterError(myFacility, 0x0015, myErrors[0x0015])
	}
	if value.name == "mainserver" {
		return performCommand(int16(code & 0x0000FFFF))
	}
	return bettererror.NewBetterError(myFacility, 0x00016, myErrors[0x0016])
}

func performCommand(code int16) error {
	return nil
}
