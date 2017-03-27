package main

import (
	"bettererror"
	"fmt"
)

type app struct {
	name     string
	commands map[string]int16
}

type installedCommands struct {
	installedApps map[int16]app
}

func (ic *installedCommands) registerApp(appFacility int16, appName string, comms map[string]int16) error {
	_, existed := ic.installedApps[appFacility]
	if existed {
		return bettererror.NewBetterError(myFacility, 0x0014, fmt.Sprintf(myErrors[0x0014], appFacility, appName))
	}
	ic.installedApps[appFacility] = app{name: appName, commands: comms}
	return nil
}

func (ic *installedCommands) deregisterApp(appFacility int16) {
	delete(ic.installedApps, appFacility)
}

func (ic *installedCommands) flattenCommands() map[string]int32 {
	var res = make(map[string]int32)
	for key, value := range ic.installedApps {
		for comm, number := range value.commands {
			res[comm] = int32(((uint32(key) << 16) & 0xFFFF0000) | (uint32(number) & 0xFFFF))
		}
	}
	return res
}
