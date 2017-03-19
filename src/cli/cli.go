package main

import (
	"args"
	"bettererror"
	"fmt"
	"messages"
	"net"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"
)

const myFacility uint16 = 0x1003

var ipaddr []net.IP
var port uint16
var myErrors = map[uint16]string{
	0x0001: "Couldn't parse/lookup IP: ",
	0x0002: "Port argument is not a valid number",
	0x0003: "Cannot use a Reserved Portnumber",
	0x0004: "Couldn't reach server: ",
	0x0005: "Write failed: ",
	0x0006: "Connection not allowed.",
	0x0007: "Read failed: ",
}

func init() {
	bettererror.RegisterFacility(myFacility, "cliapp")
}

func getHost(vs ...string) error {
	var err error
	ipaddr, err = net.LookupIP(vs[0])
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0001, myErrors[0x0001]+err.Error())
	}
	return nil
}

func setPort(vs ...string) error {
	i, err := strconv.Atoi(vs[0])
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0002, myErrors[0x0002])
	}
	if i < 1024 {
		return bettererror.NewBetterError(myFacility, 0x0003, myErrors[0x0003])
	}
	port = uint16(i)
	return nil
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

func memsetRepeat(a []byte, v byte) {
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
	memsetRepeat(data, 0)
	_, err = conn.Read(data)
	if err != nil {
		return bettererror.NewBetterError(myFacility, 0x0007, myErrors[0x0007]+err.Error())
	}
	var initResponse = new(messages.InitResponse)
	err = proto.Unmarshal(data, initResponse)
	if !initResponse.Allowed {
		return bettererror.NewBetterError(myFacility, 0x0006, myErrors[0x0006])
	}
	var b = make([]byte, 1)
	for true {
		os.Stdin.Read(b)
		fmt.Print(string(b))
	}
	return nil
}
