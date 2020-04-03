package types

import (
	"git.uangteman.com/workbench/utils/crypto"
	"strconv"
	"sync"
	"time"
)

type TokenCache struct {
	sync.Mutex
	TokenMap map[string] *AuthToken
}

type AuthToken struct {
	Token string
	TokenTimeStamp time.Time
	Key string
}

type loginResponse struct {
	Success		bool	`json:"success"`
	Message		string		`json:"message"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type signatureKeyResponse struct {
	Success		bool		`json:"success"`
	Message		string		`json:"message"`
	Data		string		`json:"data"`
}

type authRequestDTO struct {
	Username		string		`json:"username"`
	ClearPassword	string		`json:"-"`
	Password		string		`json:"password"`
	UserType		int			`json:"user_type"`
	IP				string		`json:"ip"`
	Device			string		`json:"device"`
	OS				string		`json:"os"`
	Browser			string		`json:"browser"`
	DeviceID		string		`json:"device_id"`
	Datetime		string		`json:"datetime"`
	Signature		string		`json:"signature"`
}

//func (a *authRequestDTO) GeneratePasswordSignature(keyConfig config.KeyConfig) {
//	a.GeneratePassword(keyConfig.PasswordKey)
//	a.GenerateSignature(keyConfig.SignatureKey)
//}

func (a *authRequestDTO) GeneratePassword(passKey string) {
	a.Password = crypto.EncodeSHA256HMAC(passKey, a.ClearPassword, strconv.Itoa(a.UserType))
}

func (a *authRequestDTO) GenerateSignature(signatureKey string) {
	a.Signature = crypto.EncodeSHA256HMAC(signatureKey, a.Username, a.Password, strconv.Itoa(a.UserType), a.IP, a.Device, a.OS, a.Browser, a.DeviceID, a.Datetime)
}

//func (authToken *AuthToken) Login(authConfig config.AuthUserConfig, keyConfig config.KeyConfig, integrationConfig *config.IntegrationConfig) (error) {
//	netTransport := &http.Transport{
//		DialContext: (&net.Dialer{Timeout: time.Duration(integrationConfig.HttpDialTimeoutSeconds) * time.Second}).DialContext,
//		TLSHandshakeTimeout: time.Duration(integrationConfig.HttpDialTimeoutSeconds) * time.Second,
//	}
//
//	netClient := &http.Client{Timeout: time.Duration(integrationConfig.HttpRequestTimeoutSeconds) * time.Second, Transport: netTransport}
//
//	exHttpServiceConfig := integrationConfig.Externals.Http["idm"]
//	idmLoginUrl := url.URL{
//		Scheme:     exHttpServiceConfig.Scheme,
//		Host:       exHttpServiceConfig.Host,
//		Path:       exHttpServiceConfig.Endpoints["login"],
//	}
//
//	ts := time.Now().Format("2006-01-02 15:04:05")
//
//	hostname,err := os.Hostname()
//	if hostname == "" {
//		hostname = "localhost"
//	}
//
//	authReq := authRequestDTO {
//		Username		: authConfig.User,
//		ClearPassword	: authConfig.Password,
//		UserType		: authConfig.UserType,
//		IP				: "ut-loanservice",
//		Device			: "Microservice",
//		OS				: "linux",
//		Browser			: "apps",
//		DeviceID		: hostname,
//		Datetime		: ts,
//	}
//
//	authReq.GeneratePasswordSignature(keyConfig)
//
//	var buf bytes.Buffer
//	_ = json.NewEncoder(&buf).Encode(authReq)
//
//	request, err := http.NewRequest("POST", idmLoginUrl.String(), &buf)
//	if err != nil {
//		return fmt.Errorf("failed to build request to IDM: %s", err)
//	}
//
//	resp, err := netClient.Do(request)
//	if err != nil {
//		return fmt.Errorf("failed to contact IDM: %s", err)
//	}
//
//	var response loginResponse
//	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
//		return fmt.Errorf("failed to decode IDM response: %s", err)
//	}
//
//	if response.Data.Token == "" {
//		return fmt.Errorf("failed to get token")
//	}
//
//	//TODO do not hardcode
//	//authToken.Token = response.Data.Token
//	authToken.TokenTimeStamp = time.Now()
//	authToken.Token = "token-repay"
//	authToken.Key = "productT34mUTis1"
//
//	return nil
//}

//func (authToken *AuthToken) GetSignatureKey(keyConfig config.KeyConfig, integrationConfig *config.IntegrationConfig) (error) {
//	netTransport := &http.Transport{
//		DialContext: (&net.Dialer{Timeout: time.Duration(integrationConfig.HttpDialTimeoutSeconds) * time.Second}).DialContext,
//		TLSHandshakeTimeout: time.Duration(integrationConfig.HttpDialTimeoutSeconds) * time.Second,
//	}
//
//	netClient := &http.Client{Timeout: time.Duration(integrationConfig.HttpRequestTimeoutSeconds) * time.Second, Transport: netTransport}
//
//	exHttpServiceConfig := integrationConfig.Externals.Http["idm"]
//	idmGetSignatureUrl := url.URL{
//		Scheme:     exHttpServiceConfig.Scheme,
//		Host:       exHttpServiceConfig.Host,
//		Path:       exHttpServiceConfig.Endpoints["get_signature_key"],
//	}
//
//	request, err := http.NewRequest("GET", idmGetSignatureUrl.String(), nil)
//	if err != nil {
//		return fmt.Errorf("failed to build request get signature key to IDM: %s", err)
//	}
//
//	datetime := authToken.TokenTimeStamp.Format("2006-01-02 15:04:05")
//	token := authToken.Token
//	signature := crypto.EncodeSHA256HMAC(keyConfig.SignatureKey, token, datetime)
//
//	request.Header.Add("datetime", datetime)
//	request.Header.Add("Authorization", fmt.Sprintf("token=%s", token))
//	request.Header.Add("signature", signature)
//
//	resp, err := netClient.Do(request)
//	if err != nil {
//		return fmt.Errorf("failed to contact IDM: %s", err)
//	}
//
//	var response signatureKeyResponse
//	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
//		return fmt.Errorf("failed to decode IDM response: %s", err)
//	}
//
//	authToken.Key = response.Data
//	return nil
//}
