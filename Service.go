package go_tiktok

import (
	"crypto/sha256"
	"fmt"
	"github.com/gofrs/uuid"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	go_token "github.com/leapforce-libraries/go_oauth2/token"
	"github.com/leapforce-libraries/go_oauth2/tokensource"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiName            string = "TikTok"
	apiUrl             string = "https://open.tiktokapis.com/v2"
	authUrl            string = "https://www.tiktok.com/v2/auth/authorize"
	tokenUrl           string = "https://open.tiktokapis.com/v2/oauth/token/"
	refreshTokenUrl    string = "https://open.tiktokapis.com/v2/oauth/token/"
	tokenHttpMethod    string = http.MethodPost
	defaultRedirectUrl string = "http://localhost:8080/oauth/redirect"
	dateTimeLayout     string = "2006-01-02T15:04:05Z"
)

type Service struct {
	clientKey     string
	clientSecret  string
	oAuth2Service *oauth2.Service
	errorResponse *ErrorResponse
	codeVerifiers map[string]string
	isDesktop     bool
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
	refreshTokenUrl_ := refreshTokenUrl
	oauth2ServiceConfig := oauth2.ServiceConfig{
		ClientId:        cfg.ClientKey,
		ClientSecret:    cfg.ClientSecret,
		RedirectUrl:     redirectUrl,
		AuthUrl:         authUrl,
		TokenUrl:        tokenUrl,
		RefreshTokenUrl: &refreshTokenUrl_,
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
		codeVerifiers: make(map[string]string),
		isDesktop:     strings.HasPrefix(redirectUrl, "http://localhost"),
	}, nil
}

func (service *Service) SetTokenSource(tokenSource tokensource.TokenSource) {
	service.oAuth2Service.SetTokenSource(tokenSource)
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
	url := service.oAuth2Service.AuthorizeUrl(&scope, nil, nil, &state)

	if service.isDesktop {
		codeVerifier := generateRandomString(50)

		service.codeVerifiers[state] = codeVerifier

		h := sha256.New()
		h.Write([]byte(codeVerifier))

		bs := h.Sum(nil)

		url += fmt.Sprintf("&code_challenge=%x&code_challenge_method=S256", bs)
	}

	return url
}

func (service *Service) ValidateToken() (*go_token.Token, *errortools.Error) {
	return service.oAuth2Service.ValidateToken()
}

func (service *Service) GetTokenFromCode(r *http.Request) *errortools.Error {
	extraData := url.Values{}

	if service.isDesktop {
		state := r.URL.Query().Get("state")
		if state != "" {
			codeVerifier, ok := service.codeVerifiers[state]
			if ok {
				extraData.Set("code_verifier", codeVerifier)
				delete(service.codeVerifiers, state)
			}
		}
	}

	return service.oAuth2Service.GetTokenFromCodeWithData(r, &extraData, nil)
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

func generateRandomString(length int) string {
	var sb strings.Builder
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
	charactersLength := len(characters)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		sb.WriteByte(characters[rand.Intn(charactersLength)])
	}

	return sb.String()
}
