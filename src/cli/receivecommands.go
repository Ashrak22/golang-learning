package main

import (
	"bettererror"
	"fmt"
	"messages"

	"github.com/golang/protobuf/proto"
)

func unmarshaller(buffer []byte) error {
	switch internalState {
	case stateInit:
		var initResponse = new(messages.InitResponse)
		err := proto.Unmarshal(buffer, initResponse)
		if err != nil {
			fmt.Println(err.Error())
		} else if !initResponse.Allowed {
			return bettererror.NewBetterError(myFacility, 0x0006, myErrors[0x0006])
		}
		internalState = stateWaitingForResponse
	case stateIdle, stateWaitingForResponse:
		var command = new(messages.MsgType)
		err := proto.Unmarshal(buffer, command)
		if err != nil {
			return bettererror.NewBetterError(myFacility, 0x0012, fmt.Sprintf(myErrors[0x0012], err.Error()))
		}
		msgType = command.Msgtype
		internalState = stateTypeReceived
	case stateTypeReceived:
		//var command proto.Message
		switch msgType {
		case messages.MsgType_CommandPush:
			receiveCommands(buffer)
			break
		case messages.MsgType_CommandResult:
			command := new(messages.CommandResult)
			err := proto.Unmarshal(buffer, command)
			if err != nil {
				return bettererror.NewBetterError(myFacility, 0x0012, fmt.Sprintf(myErrors[0x0012], err.Error()))
			}
			if command.CommandResult != 0 {
				fmt.Println(command.DisplayText)
			}
		}
		internalState = stateIdle
	}
	return nil
}

func receiveCommands(buffer []byte) error {
	var msg = new(messages.CommandPush)
	err := proto.Unmarshal(buffer, msg)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0012, fmt.Sprintf(myErrors[0x0012], err.Error()))
	}
	commands = msg.Commands
	return nil
}
