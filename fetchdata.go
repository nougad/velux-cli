package main

import "github.com/nougad/velux-cli/models"
import "github.com/go-openapi/runtime"
import "github.com/nougad/velux-cli/client"
import "github.com/nougad/velux-cli/client/operations"
import "github.com/go-openapi/swag"
import apiclient "github.com/nougad/velux-cli/client"
import httptransport "github.com/go-openapi/runtime/client"

type State struct {
	HomeId         string
	BridgeId       string
	Api            *client.VeluxActiveWithNetatmo
	Auth           runtime.ClientAuthInfoWriter
	NameForRoom    map[string]string
	RoomForName    map[string]string
	RoomForModule  map[string]string
	ModulesForRoom map[string][]string
	NameForModule  map[string]string
	ModuleForName  map[string]string
	ModuleStatus   map[string]*models.ModuleStatus
	RoomStatus     map[string]*models.RoomStatus
}

func fetchData(tokenFile string) *State {
	token := refreshToken(tokenFile)

	cfg := apiclient.DefaultTransportConfig()
	t := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	t.SetDebug(true)
	client := apiclient.New(t, nil)

	var state = &State{
		Api:            client,
		BridgeId:       BridgeId,
		Auth:           httptransport.BearerToken(token.AccessToken),
		NameForRoom:    make(map[string]string),
		RoomForName:    make(map[string]string),
		RoomForModule:  make(map[string]string),
		ModulesForRoom: make(map[string][]string),
		NameForModule:  make(map[string]string),
		ModuleForName:  make(map[string]string),
		ModuleStatus:   make(map[string]*models.ModuleStatus),
		RoomStatus:     make(map[string]*models.RoomStatus),
	}

	r, err := state.Api.Operations.HomesData(operations.NewHomesDataParams(), state.Auth)
	if err != nil {
		panic(err)
	}

	state.HomeId = r.Payload.Body.Homes[0].ID

	for _, r := range r.Payload.Body.Homes[0].Rooms {
		state.NameForRoom[r.ID] = r.Name
		state.NameForRoom[r.Name] = r.ID
		for _, m := range r.Modules {
			state.RoomForModule[m] = r.ID
		}
		state.ModulesForRoom[r.ID] = r.Modules
	}

	for _, m := range r.Payload.Body.Homes[0].Modules {
		state.NameForModule[m.ID] = m.Name
		state.ModuleForName[m.Name] = m.ID
	}

	param := operations.NewHomeStatusParams()
	param.WithBody(operations.HomeStatusBody{HomeID: swag.String(state.HomeId)})

	r2, err := state.Api.Operations.HomeStatus(param, state.Auth)
	if err != nil {
		panic(err)
	}

	for _, m := range r2.Payload.Body.Home.Modules {
		state.ModuleStatus[m.ID] = m
	}

	for _, r := range r2.Payload.Body.Home.Rooms {
		state.RoomStatus[r.ID] = r
	}

	return state
}
