package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"makerchecker-api/models"
	"net/http"

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
        "path": fmt.Sprintf("/%v/%v", apiRoute, id),
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
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        panic(err)
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        panic(err)
    }

    return response.StatusCode, jsonObject
}

func GetListofUsersWithRolesWithMicroservice(checkerRoles []string) (int, []map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        panic(err)
    }

    client := lambda.NewFromConfig(cfg)

    bodyJSON := map[string]interface{}{
        "roles": checkerRoles,
    }

    body, err := json.Marshal(bodyJSON)

    event := map[string]interface{}{
        "httpMethod": "POST",
        "path": "/users/accounts/with-roles",
        "body": string(body),
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

    var response models.Response

    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        panic(err)
    }

    var jsonObject []map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        panic(err)
    }

    return response.StatusCode, jsonObject
}

func UpdateWithMicroservice(lambdaFn string, apiRoute string, bodyJSON map[string]interface{}) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        panic(err)
    }

    id := bodyJSON["id"]

    client := lambda.NewFromConfig(cfg)

    body, err := json.Marshal(bodyJSON)

    event := map[string]interface{}{
        "httpMethod": "PUT",
        "path": fmt.Sprintf("/%v/%v", apiRoute, id),
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
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        panic(err)
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        panic(err)
    }

    return response.StatusCode, jsonObject
}

func CreateWithMicroservice(lambdaFn string, apiRoute string, bodyJSON map[string]interface{}) (int, map[string]interface{}) {
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
        "path": fmt.Sprintf("/%v", apiRoute),
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
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        panic(err)
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        panic(err)
    }

    return response.StatusCode, jsonObject
}

func DeleteFromMicroserviceById(lambdaFn string, apiRoute string, id string) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        panic(err)
    }

    client := lambda.NewFromConfig(cfg)

    event := map[string]interface{}{
        "httpMethod": "DELETE",
        "path": fmt.Sprintf("/%v/%v", apiRoute, id),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        return 500, map[string]interface{}{"data": "Unable to Marshal data"}
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        return 500, map[string]interface{}{"data": err.Error()}
    }

    var response models.Response
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        panic(err)
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        panic(err)
    }

    return response.StatusCode, jsonObject
}
