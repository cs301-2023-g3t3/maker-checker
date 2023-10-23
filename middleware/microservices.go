package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

var apiGateway = os.Getenv("APIGATEWAY")
var region = os.Getenv("REGION")
var stage = os.Getenv("STAGE")
var uri = fmt.Sprintf("https://%s.execute-api.%s.amazonaws.com/%s/api/v1/", apiGateway, region, stage)

func GetAllPoints() {
    routeName := fmt.Sprintf("%s/%s", uri, "points")

    res, err := http.Get(routeName)
    if err != nil {
        fmt.Print(err.Error())
        return
    }

    data, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(data))
}

func GetFromMicroserviceById(lambdaFn string, apiRoute string, id string) []byte {
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
    fmt.Println(string(res.Payload))

    return res.Payload
}

func GetUsers() []byte {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("REGION")))
    if err != nil {
        panic(err)
    }

    client := lambda.NewFromConfig(cfg)

    event := map[string]interface{}{
        "httpMethod": "GET",
        "path": fmt.Sprintf("api/v1/users"),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        panic(err)
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String("user-storage-api"),
        Payload: eventJSON,
    })

    if err != nil {
        panic(err)
    }
    fmt.Println(string(res.Payload))

    return res.Payload
}
