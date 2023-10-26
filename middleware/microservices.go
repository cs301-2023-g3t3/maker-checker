package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"makerchecker-api/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

var region = "ap-southeast-1"

func GetFromMicroserviceById(lambdaFn string, apiRoute string, id string) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        panic(err)
    }

    client := lambda.NewFromConfig(cfg)

    event := map[string]interface{}{
        "httpMethod": "GET",
        "path": fmt.Sprintf("/api/v1/%v/%v", apiRoute, id),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        panic(err)
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        panic(err)
    }

    var response models.Response
    
    json.Unmarshal(res.Payload, &response)

    var responseBody map[string]interface{}
    json.Unmarshal([]byte(response.Body), &responseBody)

    return response.StatusCode, responseBody
}

func UpdateMicroserviceById(lambdaFn string, apiRoute string, bodyJSON map[string]interface{}) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        panic(err)
    }

    id := bodyJSON["id"]

    client := lambda.NewFromConfig(cfg)

    body, err := json.Marshal(bodyJSON)

    event := map[string]interface{}{
        "httpMethod": "PUT",
        "path": fmt.Sprintf("/api/v1/%v/%v", apiRoute, id),
        "body": string(body),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        panic(err)
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        panic(err)
    }

    var response models.Response
    
    json.Unmarshal(res.Payload, &response)

    var responseBody map[string]interface{}
    json.Unmarshal([]byte(response.Body), &responseBody)

    return response.StatusCode, responseBody
}

func CreateServiceWithMicroservice(lambdaFn string, apiRoute string, bodyJSON map[string]interface{}) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        panic(err)
    }

    if _, found := bodyJSON["id"]; found {
        delete(bodyJSON, "id")
    }

    client := lambda.NewFromConfig(cfg)

    body, err := json.Marshal(bodyJSON)

    event := map[string]interface{}{
        "httpMethod": "POST",
        "path": fmt.Sprintf("/api/v1/%v", apiRoute),
        "body": string(body),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        panic(err)
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        panic(err)
    }

    var response models.Response
    
    json.Unmarshal(res.Payload, &response)

    var responseBody map[string]interface{}
    json.Unmarshal([]byte(response.Body), &responseBody)

    return response.StatusCode, responseBody
}
