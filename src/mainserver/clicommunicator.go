package main

import (
	"bettererror"
	"communicator"
	"fmt"
	"messages"
	"strconv"
	"strings"
	"time"
)

func handleCli(comm *communicator.ServerCommunicator) {
	var initResponse = &messages.InitResponse{Magic: 0xABCD, Allowed: true}
	var name string
	if strings.Contains(comm.GetRemoteAddress(), ":") {
		name = "[" + comm.GetRemoteAddress() + "]:" + strconv.Itoa(comm.GetLocalPort())
	} else {
		name = comm.GetRemoteAddress() + ":" + strconv.Itoa(comm.GetLocalPort())
	}
	climap[name] = comm
	defer removeConnection(name)
	go sendCommands("All")
	comm.Write(initResponse)
	if err := <-comm.ErrChannel; err != nil {
		fmt.Println(err.Error())
		return
	}

	for true {
		var command = new(messages.Command)

		if err := comm.Read(command); err != nil {
			if err.(*bettererror.BetterError).Code() == 0x00020007 {
				fmt.Println("Client disconnected")
				return
			}
			fmt.Println(err.Error())
			return
		}

		var resp *messages.CommandResult
		if command.Magic != 0xABCD {
			err := bettererror.NewBetterError(myFacility, 0x0009, fmt.Sprintf("%s: 0x%4X", myErrors[0x0009], command.Magic))
			resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: int32(err.Code()), DisplayText: err.Error()}
		} else {
			var err error
			fmt.Printf("Received command 0x%.8X with args '%s'\r\n", command.Command, command.Argstring)
			switch command.Command {
			case 0x0005:
				go resendCommands(name)
				err = nil
			default:
				err = nil
			}
			if err != nil {
				resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: int32(err.(*bettererror.BetterError).Code()), DisplayText: err.Error()}
			} else {
				resp = &messages.CommandResult{Magic: 0xABCD, CommandResult: 0}
			}
		}

		comm.Lock()
		respType := &messages.MsgType{Magic: 0xABCD, Msgtype: messages.MsgType_CommandResult}
		comm.Write(respType)
		comm.Write(resp)
		comm.Unlock()

		if err := <-comm.ErrChannel; err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010]+err.Error())
			fmt.Println(err.Error())
			return
		}

		if err := <-comm.ErrChannel; err != nil {
			err = bettererror.NewBetterError(myFacility, 0x0010, myErrors[0x0010]+err.Error())
			fmt.Println(err.Error())
			return
		}
	}
}

func resendCommands(conn string) error {
	sendCommands(conn)
	return nil
}

func sendCommands(conn string) {
	time.Sleep(100 * time.Millisecond)
	if conn == "All" {
		for _, value := range climap {
			var msgType = &messages.MsgType{Magic: 0xABCD, Msgtype: messages.MsgType_CommandPush}
			var msg = new(messages.CommandPush)
			msg.Magic = 0xABCD
			msg.Commands = apps.flattenCommands()

			value.Lock()
			value.Write(msgType)
			value.Write(msg)
			value.Unlock()

			if err := <-value.ErrChannel; err != nil {
				fmt.Println(err.Error())
			}
			if err := <-value.ErrChannel; err != nil {
				fmt.Println(err.Error())
			}
		}
	} else {
		var msg = new(messages.CommandPush)
		msg.Magic = 0xABCD
		msg.Commands = apps.flattenCommands()
		var msgType = &messages.MsgType{Magic: 0xABCD, Msgtype: messages.MsgType_CommandPush}

		climap[conn].Lock()
		climap[conn].Write(msgType)
		climap[conn].Write(msg)
		climap[conn].Unlock()

		if err := <-climap[conn].ErrChannel; err != nil {
			fmt.Println(err.Error())
		}

		if err := <-climap[conn].ErrChannel; err != nil {
			fmt.Println(err.Error())
		}
	}

}

func removeConnection(conn string) {
	climap[conn].Lock()
	climap[conn].Close()
	climap[conn].Unlock()
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
