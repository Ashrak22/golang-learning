package main

import (
	"bettererror"
	"fmt"
	"messages"
	"net"
)

func receiveCommands() {
	sock, err := net.Listen("tcp", "0.0.0.0:40000")
	if err != nil {
		panic(bettererror.NewBetterError(myFacility, 0x0012, fmt.Sprintf(myErrors[0x0012], err.Error())))
	}
	for true {
		conn, _ := sock.Accept()
		var msg = new(messages.CommandPush)
		if err = messages.ReadMessage(conn.(*net.TCPConn), msg, compress); err != nil {
			conn.Close()
			fmt.Println(err.Error())
			continue
		}
		commands = msg.Commands
		conn.Close()
	}
}
