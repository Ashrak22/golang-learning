package main

import (
	"bettererror"
	"bytes"
	"compress/gzip"
	"fmt"
	"messages"

	"github.com/golang/protobuf/proto"
)

const myFacility uint16 = 0x1000

func init() {
	bettererror.RegisterFacility(myFacility, "testgoApp")
}

func main() {
	var msg = &messages.CommandPush{Magic: 0xABCD, Commands: messages.Commands}
	var b bytes.Buffer
	marshalled, _ := proto.Marshal(msg)
	fmt.Println(len(marshalled))

	var gzipper = gzip.NewWriter(&b)
	gzipper.Write(marshalled)
	gzipper.Close()

	fmt.Println(b.Len())

	ungzipper, _ := gzip.NewReader(&b)
	var unzipped = make([]byte, 100*1024)
	length, _ := ungzipper.Read(unzipped)
	fmt.Println(len(unzipped), length)
	var msgun = new(messages.CommandPush)
	err := proto.Unmarshal(unzipped[:length], msgun)
	if err != nil {
		fmt.Println(err.Error())
	}
	for key, value := range msgun.Commands {
		fmt.Println(key, value)
	}
}
