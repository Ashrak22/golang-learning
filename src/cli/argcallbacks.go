package main

import (
	"bettererror"
	"net"
	"strconv"
)

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

func setCompression(vs ...string) error {
	compress = vs[0] == "true"
	return nil
}
