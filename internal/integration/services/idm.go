package services

import (
	"github.com/yudzmaestro/test-api-microservice/internal/config"
	"git.uangteman.com/workbench/utils/http/client"
)

type IDMHttpService struct {
	config     *config.ExternalHttpServiceConfig
	configauth *config.AuthConfig
	netClient  *client.HttpClient
}

func NewIDMHttpIntegrationService(config *config.ExternalHttpServiceConfig, configAuth *config.AuthConfig, httpClient *client.HttpClient) (IDMHttpService, error) {
	ls := IDMHttpService{
		config:     config,
		configauth: configAuth,
		netClient:  httpClient,
	}
	return ls, nil
}
