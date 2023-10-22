package models

type Response struct {
    Body                string                  `json:"body"`
    Headers             map[string]interface{}  `json:"headers"`
    MultiValueHeaders   map[string]interface{}  `json:"multiValueHeaders"`
    StatusCode          int                     `json:"statusCode"`
}
