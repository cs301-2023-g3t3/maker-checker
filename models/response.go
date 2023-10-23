package models

type Response struct {
    StatusCode        int                    `json:"statusCode"`
    Headers           map[string]string      `json:"headers"`
    MultiValueHeaders map[string]interface{}    `json:"multiValueHeaders"`
    Body              string                 `json:"body"`
}
