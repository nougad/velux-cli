package main

import "fmt"

func PrintStatus(state *State) {
	for _, r := range state.RoomStatus {
		if r.Temperature != 0 {
			fmt.Printf(
				"%s (air quality: %d / CO2: %d / Temperature: %d / Humidity: %d / Lux: %d)\n",
				state.NameForRoom[r.Id], r.AirQuality, r.Co2, r.Temperature/10.0, r.Humidity, r.Lux)
		} else {
			fmt.Printf("%s\n", state.NameForRoom[r.Id])
		}

		for _, m := range state.ModulesForRoom[r.Id] {
			if state.ModuleStatus[m].Type_ == "NXO" {
				fmt.Printf("  - %d %s\n", state.ModuleStatus[m].CurrentPosition, state.NameForModule[m])
			} else {
				fmt.Printf("  - %s: battery: %d%% rf strength: %d\n", state.NameForModule[m], state.ModuleStatus[m].BatteryPercent, state.ModuleStatus[m].RfStrength)
			}
		}
	}
}
