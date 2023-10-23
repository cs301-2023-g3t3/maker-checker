package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"makerchecker/models"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func GetFromMicroserviceById(lambdaFn string, apiRoute string, id string) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
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
