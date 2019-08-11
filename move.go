package main

import "fmt"
import "encoding/json"
import sw "./go-client"

func Move(state *State, shutters []string, position int32) {
	fmt.Printf("Moving shutters: %+v to %+v\n", shutters, position)
	if len(shutters) == 0 {
		return
	}

	var updates []sw.ModulePercentage
	for _, x := range shutters {
		m := sw.ModulePercentage{
			Bridge:         state.BridgeId,
			Id:             state.ModuleForName[x],
			TargetPosition: position,
		}
		updates = append(updates, m)
	}

	param := sw.SetState{
		Home: &sw.SetStateHome{
			Id:      state.HomeId,
			Modules: updates,
		},
	}

	fmt.Printf("> request: %+v\n", param)
	fmt.Printf("> request: %+v\n", param.Home)
	fmt.Printf("> request: %+v\n", param.Home.Modules)
	j, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	fmt.Printf("> request: %+v\n", string(j))

	response, _, err := state.Api.SetState(state.Auth, param)
	if err != nil {
		panic(err)
	}

	j, err = json.Marshal(response)
	if err != nil {
		panic(err)
	}
	fmt.Printf("> response: %+v\n", string(j))
}
