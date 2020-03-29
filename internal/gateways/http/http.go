package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	//"net"
	//"net/url"
	"os"
	"strings"
	//"time"
	"errors"

	"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	//fasthttptransport "github.com/l-vitaly/go-kit/transport/fasthttp"

	"github.com/yudzmaestro/test-api-microservice/internal/config"
	"github.com/yudzmaestro/test-api-microservice/internal/handlers"
	"github.com/yudzmaestro/test-api-microservice/internal/models"
	"github.com/yudzmaestro/test-api-microservice/pkg/errs"

	"github.com/spartaut/utils/http/helper"
	"github.com/spartaut/utils/http/router"
	"github.com/spartaut/utils/middlewares"

	//imiddlewares "github.com/spartaut/ut_promo/internal/middlewares"

	"net/http"
	"reflect"

	"io/ioutil"
)

type baseHttpResponse struct {
	RequestID string      `json:"request_id"`
	Code      int         `json:"code"`
	Error     string      `json:"error, omitempty"`
}

func contextInjector(ctx context.Context, r *http.Request) context.Context {
	// get token
	var auth = r.Header.Get("Authorization")

	// validate token
	var tokens = strings.Split(auth, "=")

	var token = strings.TrimSpace(tokens[1])

	// get datetime
	var datetime = r.Header.Get("datetime")

	// get signature
	var signature = r.Header.Get("signature")

	authObject := models.AuthObject{
		Token:     token,
		Datetime:  datetime,
		Signature: signature,
	}
	ctx = context.WithValue(ctx, "auth-object", authObject)

	return ctx
}

func NewHTTPServer(endpoints handlers.Endpoints, logger log.Logger, integrationConfig *config.IntegrationConfig) *router.Router {
	r := router.NewRouter()

	r.Use(middlewares.HttpRequestIDInjectorMiddleware)
	//r.Use(utmiddlewares.HttpRequestIDInjectorMiddleware)

	//options := []httptransport.ServerOption{
	//	httptransport.ServerErrorEncoder(errorEncoder),
	//	httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	//	httptransport.ServerBefore(helper.PopulateRequestContext),
	//}
	//options := []httptransport.ServerOption{}

	version := os.Getenv("VERSION")
	if version == "" {
		version = "pong"
	}

	//netTransport := &http.Transport{
	//	DialContext: (&net.Dialer{Timeout: time.Duration(integrationConfig.HttpDialTimeoutSeconds) * time.Second}).DialContext,
	//	TLSHandshakeTimeout: time.Duration(integrationConfig.HttpDialTimeoutSeconds) * time.Second,
	//}

	//netClient := &http.Client{Timeout: time.Duration(integrationConfig.HttpRequestTimeoutSeconds) * time.Second, Transport: netTransport}

	//idmAuthorizeUrl := url.URL{
	//	Scheme:     integrationConfig.Externals.Http["idm"].Scheme,
	//	Host:       integrationConfig.Externals.Http["idm"].Host,
	//	Path:       integrationConfig.Externals.Http["idm"].Endpoints["authorize"],
	//}

	httpLoggingMiddleware := middlewares.MakeHttpTransportLoggingMiddleware(logger)
	//httpLoggingWithAuthMiddlewareFunc := middlewares.MakeHttpTransportLoggingWithAuthMiddleware(idmAuthorizeUrl.String(), netClient, integrationConfig.UTKey, logger)

	//systemSecure := imiddlewares.MakeSecureHttpRequestMiddleware(tokenCache.TokenMap["system"])

	s := r.Subroute("/loan")
	s.Methods("GET", "POST").Handler("/ping", httpLoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(version))
	})))

	//s.Methods("POST").Handler("/register-loan", httpLoggingWithAuthMiddlewareFunc(
	//	httptransport.NewServer(
	//		endpoints.RegisterLoanEndpoint,
	//		decodeRegisterLoanRequest,
	//		commonEncodeResponse,
	//		options...,
	//	), "loanPostRegisterLoan",""))
	//
	//s.Methods("GET").Handler( "/app_loan_id/:id", httpLoggingWithAuthMiddlewareFunc(
	//	httptransport.NewServer(
	//		endpoints.GetLoanByAppLoanIDEndpoint,
	//		decodeGetLoanByAppLoanIDRequest,
	//		commonEncodeResponse,
	//		options...,
	//	), "loanGetByAppLoanID",""))

	return r
}

func WithHttpMiddlewares(httpHandler http.Handler, middlewares ...middlewares.MiddlewareFunc) http.Handler {
	for _, mw := range middlewares {
		httpHandler = mw(httpHandler)
	}

	return httpHandler
}

func makeCommonPostRequestDecoder(model interface{}) httptransport.DecodeRequestFunc {
	req := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			return nil, err
		}

		return req, nil
	}
}

func commonXMLDecodeRequest(model interface{}) httptransport.DecodeRequestFunc {
	req := reflect.New(reflect.TypeOf(model).Elem()).Interface()

	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		bodyAsByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = xml.Unmarshal(bodyAsByte, &req)
		if err != nil {
			return nil, err
		}

		return req, nil
	}
}

func commonXMLEncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "text/xml")
	return xml.NewEncoder(w).Encode(response)
}

func commonEncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func commonStringEncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	msg := fmt.Sprintf("%v", response)
	_, err := w.Write([]byte(msg))
	return err
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	code := err2code(err)
	w.WriteHeader(code)
	reqid, ok := helper.ReqIDFromContext(ctx)
	if !ok {
		reqid = ""
	}
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error(), ReqID: reqid, Code: code})
}

func err2code(err error) int {
	if errors.Is(err, errs.ErrBadRequest) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	ReqID string `json:"request_id"`
	Code  int    `json:"code"`
	Error string `json:"error"`
}
