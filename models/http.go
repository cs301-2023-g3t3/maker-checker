package models

type HttpResponse struct {
	Code        int                    `json:"code" example:"400"`
	Message     string                 `json:"message" example:"status bad request"`
    Data        map[string]interface{} `json:"data"`
}
