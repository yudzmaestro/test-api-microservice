package services

import (
	"context"

	"github.com/spartaut/utils/http/client"
	"github.com/yudzmaestro/test-api-microservice/internal/config"
	"github.com/yudzmaestro/test-api-microservice/internal/integration/dto"
	"net/url"
	"github.com/spartaut/utils/http/helper"
	"encoding/json"
	"fmt"
	"github.com/spartaut/utils/commons"
	"errors"
)

type NotificationService interface{
	SendNotif(ctx context.Context, notifRequest *dto.NotificationRequestDTO) (error)
}

type notificationHttpService struct {
	config *config.ExternalHttpServiceConfig
	netClient *client.HttpClient
}

func NewNotificationHttpIntegrationService(config *config.ExternalHttpServiceConfig, httpClient *client.HttpClient) (NotificationService, error) {
	ns := &notificationHttpService{
		config:    config,
		netClient: httpClient,
	}

	return ns, nil
}

func (rs *notificationHttpService) SendNotif(ctx context.Context, notifRequest *dto.NotificationRequestDTO) (error) {

	auth,err := helper.AuthObjectFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to call repayment service: auth object is not available in context. err :%s", err)
	}

	ep := url.URL{
		Scheme:     rs.config.Scheme,
		Host:       rs.config.Host,
		Path:       rs.config.Endpoints[config.URL_NOTIFICATION_SEND_NOTIF],
	}

	notifRequest.Datetime = auth.Datetime
	notifRequest.Signature = auth.Signature

	dataRequest,_ := json.Marshal(notifRequest)
	httpCall := commons.HttpCall{
		Method		: commons.HTTP_CALL_METHOD_POST,
		URL			: ep.String(),
		DataRequest	: dataRequest,
		ContentType	: commons.HTTP_CALL_CONTENT_JSON,
		Headers		: auth.ToHeadersMap(),
	}

	dataResponse ,err := httpCall.SendRequest()
	if err != nil {
		return err
	}

	response := dto.NotificationResponseDTO{}
	err = json.Unmarshal(dataResponse, &response)
	if err != nil {
		return err
	}

	if response.Code != 0 {
		return errors.New(response.Message)
	}

	return nil
}