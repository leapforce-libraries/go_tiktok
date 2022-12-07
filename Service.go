package go_tiktok

import (
	"fmt"
	"github.com/gofrs/uuid"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	go_token "github.com/leapforce-libraries/go_oauth2/token"
	"github.com/leapforce-libraries/go_oauth2/tokensource"
	"net/http"
	"time"
)

const (
	apiName            string = "TikTok"
	apiUrl             string = "https://open.tiktokapis.com/v2"
	authUrl            string = "https://www.tiktok.com/auth/authorize"
	tokenUrl           string = "https://open-api.tiktok.com/oauth/access_token"
	tokenHttpMethod    string = http.MethodPost
	defaultRedirectUrl string = "http://localhost:8080/oauth/redirect"
	dateTimeLayout     string = "2006-01-02T15:04:05Z"
)

type Service struct {
	clientKey     string
	clientSecret  string
	oAuth2Service *oauth2.Service
	errorResponse *ErrorResponse
}

type ServiceConfig struct {
	ClientKey     string
	ClientSecret  string
	TokenSource   tokensource.TokenSource
	RedirectUrl   *string
	RefreshMargin *time.Duration
}

func NewService(cfg *ServiceConfig) (*Service, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if cfg.ClientKey == "" {
		return nil, errortools.ErrorMessage("ClientKey not provided")
	}

	redirectUrl := defaultRedirectUrl
	if cfg.RedirectUrl != nil {
		redirectUrl = *cfg.RedirectUrl
	}

	clientIdName := "client_key"
	oauth2ServiceConfig := oauth2.ServiceConfig{
		ClientId:        cfg.ClientKey,
		ClientSecret:    cfg.ClientSecret,
		RedirectUrl:     redirectUrl,
		AuthUrl:         authUrl,
		TokenUrl:        tokenUrl,
		RefreshMargin:   cfg.RefreshMargin,
		TokenHttpMethod: tokenHttpMethod,
		TokenSource:     cfg.TokenSource,
		ClientIdName:    &clientIdName,
	}
	oauth2Service, e := oauth2.NewService(&oauth2ServiceConfig)
	if e != nil {
		return nil, e
	}

	return &Service{
		clientKey:     cfg.ClientKey,
		oAuth2Service: oauth2Service,
	}, nil
}

func (service *Service) HttpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	var request *http.Request
	var response *http.Response
	var e *errortools.Error

	// add error model
	service.errorResponse = &ErrorResponse{}
	requestConfig.ErrorModel = service.errorResponse

	request, response, e = service.oAuth2Service.HttpRequest(requestConfig)
	if e != nil {
		if service.errorResponse.Data.Description != "" {
			e.SetMessage(service.errorResponse.Data.Description)
		}
	}

	return request, response, e
}

func (service *Service) AuthorizeUrl(scope string) string {
	g := uuid.NewGen()
	guid, _ := g.NewV1()
	state := guid.String()
	return service.oAuth2Service.AuthorizeUrl(scope, nil, nil, &state)
}

func (service *Service) ValidateToken() (*go_token.Token, *errortools.Error) {
	return service.oAuth2Service.ValidateToken()
}

func (service *Service) GetTokenFromCode(r *http.Request) *errortools.Error {
	return service.oAuth2Service.GetTokenFromCode(r, nil)
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiUrl, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.clientKey
}

func (service *Service) ApiCallCount() int64 {
	return service.oAuth2Service.ApiCallCount()
}

func (service *Service) ApiReset() {
	service.oAuth2Service.ApiReset()
}

func (service *Service) ErrorResponse() *ErrorResponse {
	return service.errorResponse
}
