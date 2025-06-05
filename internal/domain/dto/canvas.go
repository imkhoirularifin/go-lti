package dto

type Oauth2RedirectRequest struct {
	Code             string `query:"code"`
	State            string `query:"state"`
	Error            string `query:"error"`
	ErrorDescription string `query:"error_description"`
}

type Oauth2ExchangeResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	User        struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ExpiresIn    int    `json:"expires_in"`
	CanvasRegion string `json:"canvas_region"`
}
