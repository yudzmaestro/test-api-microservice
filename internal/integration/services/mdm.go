package services

import (
	"github.com/yudzmaestro/test-api-microservice/internal/config"
	"github.com/spartaut/utils/http/client"
)

type MDMHttpService struct {
	config			*config.ExternalHttpServiceConfig
	configauth 		*config.AuthConfig
	netClient		*client.HttpClient
}

func NewMDMHttpIntegrationService(config *config.ExternalHttpServiceConfig, configauth *config.AuthConfig, httpClient *client.HttpClient) (MDMHttpService, error)  {
	ls := MDMHttpService{
		config:     config,
		configauth: configauth,
		netClient:  httpClient,
	}

	return ls, nil
}

//Get Private Key

//Get Public Key

//Get Signature Key
