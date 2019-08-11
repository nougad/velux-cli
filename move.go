package main

import "fmt"
import "encoding/json"
import "github.com/nougad/velux-cli/models"
import "github.com/nougad/velux-cli/client/operations"
import "github.com/go-openapi/swag"

func Move(state *State, shutters []string, position int64) {
	fmt.Printf("Moving shutters: %+v to %+v\n", shutters, position)
	if len(shutters) == 0 {
		return
	}

	var updates []*models.ModulePercentage
	for _, x := range shutters {
		m := &models.ModulePercentage{
			Bridge:         state.ModuleStatus[state.ModuleForName[x]].Bridge,
			ID:             state.ModuleForName[x],
			TargetPosition: swag.Int64(position),
		}
		updates = append(updates, m)
	}

	param := operations.NewSetStateParams()
	param.WithBody(&models.SetState{
		Home: &models.SetStateHome{
			ID:      swag.String(state.HomeId),
			Modules: updates,
		},
	})

	fmt.Printf("> request: %+v\n", param.Body)
	fmt.Printf("> request: %+v\n", param.Body.Home)
	fmt.Printf("> request: %+v\n", param.Body.Home.Modules)
	j, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	fmt.Printf("> request: %+v\n", string(j))

	response, err := state.Api.Operations.SetState(param, state.Auth)
	if err != nil {
		panic(err)
	}

	j, err = json.Marshal(response)
	if err != nil {
		panic(err)
	}
	fmt.Printf("> response: %+v\n", string(j))
}
