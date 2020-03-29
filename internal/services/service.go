package services

import (
	"fmt"
	"github.com/yudzmaestro/test-api-microservice/internal/config"
	iservice "github.com/yudzmaestro/test-api-microservice/internal/integration/services"
	"github.com/spartaut/utils/datastore"
	"github.com/spartaut/utils/http/client"
	"time"
)

type Service interface {
	PromoService
}

type PromoService interface {
	//GetAccountAllLoan(ctx context.Context, req pkgDto.AllDataCustLoanAccountDTO) ([]models.DataLoanAll, error) // Get All Account Loan
	//
	//SaveNoteCustLoan(ctx context.Context, mdl models.DataNoteCustLoan) error              // Post notes
	//GetAllNoteCustLoan() ([]*models.DataNoteCustLoan, error)                              // get all notes
	//GetNoteCustLoan(ctx context.Context, NoteID string) (*models.DataNoteCustLoan, error) //get Notes By ID
	//
	//GetAccountDetailsLoan(ctx context.Context, req pkgDto.DetailDataCustLoanAccountDTO) (models.ResponDetailLoan, error)
	//GetPersonalInformationLoan(ctx context.Context, ID string) (models.PersonalInformation, error)
	//GetLoanInformationLoan(ctx context.Context, ID string) (models.LoanInformation, error)
	//GetPaymenHistoryLoan(ctx context.Context, Repayment_ID string) (*dto.AllRepaymentDataDTO, error)
	//GetDiscountistoryLoan(ctx context.Context, loan_ID string) ([]*models.DiscountHistoryDTO, error)
	//LoginIdm(ctx context.Context, DataDto pkgDto.LoginIdmDTO) ([]byte, error)
	//ForgetPassword(ctx context.Context, DataDto pkgDto.ForgetPasswordRequestDTO) ([]byte, error)
}

type promoService struct {
	idmService      iservice.IDMHttpService
	mdmService		iservice.MDMHttpService
	dataStore 		*datastore.DataStore
}

type service struct {
	PromoService
}

func NewService(config *config.Config) (Service, error) {
	dbConfig := datastore.DBConfig{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Dbname:   config.DB.DbName,
		User:     config.DB.User,
		Password: config.DB.Password,
	}

	dataStore, err := datastore.NewPostgresDatastore(dbConfig)

	if err != nil {
		return nil, err
	}

	collectionService, err := NewPromoService(dataStore, config)

	if err != nil {
		return nil, err
	}

	return collectionService, nil
}

func NewPromoService(ds *datastore.DataStore, conf *config.Config) (PromoService, error) {

	cached := conf.Server.EnableCache
	options := []datastore.StoreOptions{}
	if cached {
		options = append(options, datastore.Cache(500_000))
	}

	dialTimeout := time.Duration(conf.Integrations.HttpDialTimeoutSeconds) * time.Second
	reqTimeout := time.Duration(conf.Integrations.HttpRequestTimeoutSeconds) * time.Second
	netClient := client.NewWithTimeout(dialTimeout, reqTimeout)

	//idm
	idmIntegrationService, err := iservice.NewIDMHttpIntegrationService(conf.Integrations.Externals.Http["idm"], conf.Auth, netClient)
	if err != nil {
		return nil, fmt.Errorf("failed to init idm integration service: %s", err)
	}

	//mdm
	mdmIntegrationService, err := iservice.NewMDMHttpIntegrationService(conf.Integrations.Externals.Http["mdm"], conf.Auth, netClient)
	if err != nil {
		return nil, fmt.Errorf("failed to init mdm integration service: %s", err)
	}

	return &promoService{
		idmService:       idmIntegrationService,
		mdmService:		  mdmIntegrationService,
		dataStore:        ds,
	}, nil
}
