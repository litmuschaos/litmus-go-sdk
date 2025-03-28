package infrastructure

import models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type InfraData struct {
	Data   InfraList `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type InfraList struct {
	ListInfraDetails models.ListInfraResponse `json:"listInfras"`
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type InfraConnectionData struct {
	Data   RegisterInfra `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type RegisterInfra struct {
	RegisterInfraDetails models.RegisterInfraResponse `json:"registerInfra"`
}

type DisconnectInfraData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DisconnectInfraDetails `json:"data"`
}

type DisconnectInfraDetails struct {
	Message string `json:"deleteInfra"`
}

type ServerVersionResponse struct {
	Data   ServerVersionData `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type ServerVersionData struct {
	GetServerVersion models.ServerVersionResponse `json:"getServerVersion"`
}
