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
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

/*Error handling variables and consts*/
const myFacility uint16 = 0x1001

var myErrors = map[uint16]string{
	0x0001: "Port argument is not a valid number",
	0x0002: "Cannot use a Reserved Portnumber",
	0x0003: "Cannot execute screen",
	0x0004: "Cannot execute cli command",
	0x0005: "Listening Port can't be 0",
	0x0006: "Couldn't start listening on Port. ",
	0x0007: "Error when acceptin connection ",
	0x0008: "Error when waiting for data: ",
	0x0009: "Wrong Magic Number received",
}

var port uint16
var cli bool

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
		go runConnection(conn)
	}
	return nil
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

func runConnection(conn *net.TCPConn) {
	var buffer [8 * 1024]byte
	defer conn.Close()
	memsetRepeat(buffer[0:], 0)
	data, err := conn.Read(buffer[0:])
	if err != nil {
		if strings.Contains(err.Error(), "EOF") {
			return
		}
		err = bettererror.NewBetterError(myFacility, 0x0008, myErrors[0x0008]+err.Error())
		fmt.Println(err.Error())
		return
	}
	if data == 0 {
		return
	}
	var initmsg = new(messages.Init)
	err = proto.Unmarshal(buffer[0:], initmsg)
	if initmsg.Magic != 0xABCD {
		err := bettererror.NewBetterError(myFacility, 0x0009, myErrors[0x0009])
		fmt.Println(err.Error())
		return
	}
	if initmsg.App == "cli" {
		fmt.Println("App cli has connected from ", conn.RemoteAddr().String())
		handleCli(conn, buffer)
	}
}

func handleCli(conn *net.TCPConn, buffer [8 * 1024]byte) {
	for true {
		memsetRepeat(buffer[0:], 0)
		var initResponse = &messages.InitResponse{Magic: 0xABCD, Allowed: true}
		buffer, err := proto.Marshal(initResponse)
		_, err = conn.Write(buffer)
		memsetRepeat(buffer[0:], 0)
		data, err := conn.Read(buffer[0:])
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				return
			}
			err = bettererror.NewBetterError(myFacility, 0x0008, myErrors[0x0008]+err.Error())
			fmt.Println(err.Error())
			return
		}
		if data == 0 {
			break
		}
	}
}
