package main

import "encoding/json"
import "fmt"
import "io/ioutil"

type Status struct {
	ShutterStatus  map[string]int32
	AirQuality     map[string]int32
	Co2            map[string]int32
	Temperature    map[string]int32
	Humidity       map[string]int32
	Lux            map[string]int32
	BatteryPercent map[string]int32
	RfStrength     map[string]int32
}

func DumpJSON(state *State, outfile string) {
	var status = new(Status)
	status.AirQuality = make(map[string]int32)
	status.Co2 = make(map[string]int32)
	status.Temperature = make(map[string]int32)
	status.Humidity = make(map[string]int32)
	status.Lux = make(map[string]int32)
	status.ShutterStatus = make(map[string]int32)
	status.BatteryPercent = make(map[string]int32)
	status.RfStrength = make(map[string]int32)

	for _, r := range state.RoomStatus {
		roomName := state.NameForRoom[r.Id]
		if r.Temperature != 0 {
			status.AirQuality[roomName] = r.AirQuality
			status.Co2[roomName] = r.Co2
			status.Temperature[roomName] = r.Temperature
			status.Humidity[roomName] = r.Humidity
			status.Lux[roomName] = r.Lux
		}
	}

	for _, m := range state.ModuleStatus {
		moduleName := state.NameForModule[m.Id]
		if m.Type_ == "NXO" {
			status.ShutterStatus[moduleName] = m.TargetPosition
		} else if m.Type_ == "NXG" {
			// bridge -- noop
		} else {
			status.BatteryPercent[moduleName] = m.BatteryPercent
			status.RfStrength[moduleName] = m.RfStrength
		}
	}

	jsonOut, _ := json.Marshal(status)
	if outfile == "-" {
		fmt.Println(string(jsonOut))
	} else {
		err := ioutil.WriteFile(outfile, jsonOut, 0644)
		if err != nil {
			panic(err)
		}
	}
}
