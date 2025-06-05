package dto

type ResponseDto struct {
	Message     string `json:"message"`
	ErrorDetail string `json:"error_detail"`
	Data        any    `json:"data"`
}
