package main

import "context"
import sw "./go-client"

type State struct {
	HomeId         string
	BridgeId       string
	Api            *sw.DefaultApiService
	Auth           context.Context
	NameForRoom    map[string]string
	RoomForName    map[string]string
	RoomForModule  map[string]string
	ModulesForRoom map[string][]string
	NameForModule  map[string]string
	ModuleForName  map[string]string
	ModuleStatus   map[string]sw.ModuleStatus
	RoomStatus     map[string]sw.RoomStatus
}

func fetchData(tokenFile string) *State {
	token := refreshToken(tokenFile)

	var state = &State{
		Api:            sw.NewAPIClient(sw.NewConfiguration()).DefaultApi,
		BridgeId:       BridgeId,
		Auth:           context.WithValue(context.Background(), sw.ContextAccessToken, token.AccessToken),
		NameForRoom:    make(map[string]string),
		RoomForName:    make(map[string]string),
		RoomForModule:  make(map[string]string),
		ModulesForRoom: make(map[string][]string),
		NameForModule:  make(map[string]string),
		ModuleForName:  make(map[string]string),
		ModuleStatus:   make(map[string]sw.ModuleStatus),
		RoomStatus:     make(map[string]sw.RoomStatus),
	}

	r, _, err := state.Api.HomesData(state.Auth)
	if err != nil {
		panic(err)
	}

	state.HomeId = r.Body.Homes[0].Id

	for _, r := range r.Body.Homes[0].Rooms {
		state.NameForRoom[r.Id] = r.Name
		state.NameForRoom[r.Name] = r.Id
		for _, m := range r.Modules {
			state.RoomForModule[m] = r.Id
		}
		state.ModulesForRoom[r.Id] = r.Modules
	}

	for _, m := range r.Body.Homes[0].Modules {
		state.NameForModule[m.Id] = m.Name
		state.ModuleForName[m.Name] = m.Id
	}

	r2, _, err := state.Api.HomeStatus(state.Auth, sw.Body{HomeId: state.HomeId})
	if err != nil {
		panic(err)
	}

	for _, m := range r2.Body.Home.Modules {
		state.ModuleStatus[m.Id] = m
	}

	for _, r := range r2.Body.Home.Rooms {
		state.RoomStatus[r.Id] = r
	}

	return state
}
