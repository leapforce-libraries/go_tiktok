package go_tiktok

// ErrorResponse stores a general Microsoft Graph API error response
//
type ErrorResponse struct {
	Data struct {
		Captcha     string `json:"captcha"`
		DescUrl     string `json:"desc_url"`
		Description string `json:"description"`
		ErrorCode   int64  `json:"error_code"`
	} `json:"data"`
	Message string `json:"message"`
}
