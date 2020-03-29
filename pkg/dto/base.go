package dto

type BaseResponse struct {
	TraceID string `json:"trace_id"`
	Error string `json:"error"`
	Code int `json:"code"`
}
