package go_tiktok

import (
	"encoding/json"
)

type Token struct {
	Data struct {
		AccessToken      *string          `json:"access_token"`
		Captcha          string           `json:"captcha"`
		DescUrl          string           `json:"desc_url"`
		Description      string           `json:"description"`
		ErrorCode        int64            `json:"error_code"`
		ExpiresIn        *json.RawMessage `json:"expires_in"`
		LogId            string           `json:"log_id"`
		OpenId           string           `json:"open_id"`
		RefreshExpiresIn int64            `json:"refresh_expires_in"`
		RefreshToken     *string          `json:"refresh_token"`
		Scope            *string          `json:"scope"`
	} `json:"data"`
	Message *string `json:"message"`
}
