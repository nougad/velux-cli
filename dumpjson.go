package main

import "encoding/json"
import "fmt"
import "io/ioutil"

type Status struct {
	ShutterStatus  map[string]int64
	AirQuality     map[string]int64
	Co2            map[string]int64
	Temperature    map[string]int64
	Humidity       map[string]int64
	Lux            map[string]int64
	BatteryPercent map[string]int64
	RfStrength     map[string]int64
}

func DumpJSON(state *State, outfile string) {
	var status = new(Status)
	status.AirQuality = make(map[string]int64)
	status.Co2 = make(map[string]int64)
	status.Temperature = make(map[string]int64)
	status.Humidity = make(map[string]int64)
	status.Lux = make(map[string]int64)
	status.ShutterStatus = make(map[string]int64)
	status.BatteryPercent = make(map[string]int64)
	status.RfStrength = make(map[string]int64)

	for _, r := range state.RoomStatus {
		roomName := state.NameForRoom[r.ID]
		if r.Temperature != 0 {
			status.AirQuality[roomName] = r.AirQuality
			status.Co2[roomName] = r.Co2
			status.Temperature[roomName] = r.Temperature
			status.Humidity[roomName] = r.Humidity
			status.Lux[roomName] = r.Lux
		}
	}

	for _, m := range state.ModuleStatus {
		moduleName := state.NameForModule[m.ID]
		if m.Type == "NXO" {
			status.ShutterStatus[moduleName] = m.TargetPosition
		} else if m.Type == "NXG" {
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
