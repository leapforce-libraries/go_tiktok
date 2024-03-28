package go_tiktok

import "encoding/json"

type Token struct {
	AccessToken      string          `json:"access_token"`
	ExpiresIn        json.RawMessage `json:"expires_in"`
	OpenId           string          `json:"open_id"`
	RefreshExpiresIn int             `json:"refresh_expires_in"`
	RefreshToken     string          `json:"refresh_token"`
	Scope            string          `json:"scope"`
	TokenType        string          `json:"token_type"`
	Error            string          `json:"error"`
	ErrorDescription string          `json:"error_description"`
	LogId            string          `json:"log_id"`
}
