package go_tiktok

// ErrorResponse stores a general Microsoft Graph API error response
//
type ErrorResponse struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date            string `json:"date"`
			RequestId       string `json:"request-id"`
			ClientRequestId string `json:"client-request-id"`
		} `json:"innerError"`
	} `json:"error"`
}
