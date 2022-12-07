package go_tiktok

import (
	"encoding/json"
	errortools "github.com/leapforce-libraries/go_errortools"
	gcs "github.com/leapforce-libraries/go_googlecloudstorage"
	oauth2_token "github.com/leapforce-libraries/go_oauth2/token"
	"github.com/leapforce-libraries/go_oauth2/tokenmap"
)

type TokenMap struct {
	tm *tokenmap.TokenMap
}

func NewTokenMap(map_ *gcs.Map) (*TokenMap, *errortools.Error) {
	tm, e := tokenmap.NewTokenMap(map_)
	if e != nil {
		return nil, e
	}

	return &TokenMap{tm}, nil
}

func (m *TokenMap) Token() *oauth2_token.Token {
	return m.tm.Token()
}

func (m *TokenMap) NewToken() (*oauth2_token.Token, *errortools.Error) {
	return m.tm.NewToken()
}

func (m *TokenMap) SetToken(token *oauth2_token.Token, save bool) *errortools.Error {
	return m.tm.SetToken(token, save)
}

func (m *TokenMap) RetrieveToken() *errortools.Error {
	return m.tm.RetrieveToken()
}

func (m *TokenMap) SaveToken() *errortools.Error {
	return m.tm.SaveToken()
}

func (m *TokenMap) UnmarshalToken(b []byte) (*oauth2_token.Token, *errortools.Error) {
	var token Token

	err := json.Unmarshal(b, &token)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	if token.Message != nil {
		if *token.Message == "error" {
			return nil, errortools.ErrorMessage(token.Data.Description)
		}
	}

	token_ := oauth2_token.Token{
		AccessToken:  token.Data.AccessToken,
		Scope:        token.Data.Scope,
		ExpiresIn:    token.Data.ExpiresIn,
		RefreshToken: token.Data.RefreshToken,
	}

	return &token_, nil
}
