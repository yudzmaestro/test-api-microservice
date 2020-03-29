package handlers

type Endpoints struct {
	PromoEndpoints
}

type baseResponse struct {
	ReqID string `json:"request_id"`
	Code int `json:"code"`
	Error string `json:"error,omitempty"`
}

type baseMobileResponse struct {
	ReqID string `json:"request_id"`
	Code int `json:"code"`
	Message string `json:"message,omitempty"`
	Error string `json:"error,omitempty"`
}