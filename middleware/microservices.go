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
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    client := lambda.NewFromConfig(cfg)

    event := map[string]interface{}{
        "httpMethod": "GET",
        "path": fmt.Sprintf("/%v/%v", apiRoute, id),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    var response models.Response
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    return response.StatusCode, jsonObject
}

func GetListofUsersWithRolesWithMicroservice(checkerRoles []float64, idToken string) (int, any) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }

    client := lambda.NewFromConfig(cfg)

    bodyJSON := map[string]interface{}{
        "roles": checkerRoles,
    }

    body, err := json.Marshal(bodyJSON)

    event := map[string]interface{}{
        "httpMethod": "POST",
        "path": "/users/accounts/with-roles",
        "headers": map[string]string{
            "Content-Type": "application/json",
            "X-IDTOKEN": idToken,
        },
        "body": string(body),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String("user-storage-api"),
        Payload: eventJSON,
    })

    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }

    var response models.Response

    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }

    if response.Body == "" {
        return 404, "No available checker found"
    }

    var jsonObject []map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        return http.StatusInternalServerError, err.Error()
    }

    return response.StatusCode, jsonObject
}

func UpdateWithMicroservice(lambdaFn string, apiRoute string, bodyJSON map[string]interface{}, idToken string) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    id := bodyJSON["id"]

    client := lambda.NewFromConfig(cfg)

    body, err := json.Marshal(bodyJSON)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    event := map[string]interface{}{
        "httpMethod": "PUT",
        "path": fmt.Sprintf("/%v/%v", apiRoute, id),
        "headers": map[string]string{
            "Content-Type": "application/json",
            "X-IDTOKEN": idToken,
        },
        "body": string(body),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    var response models.Response
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": "Page not found"}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    return response.StatusCode, jsonObject
}

func CreateWithMicroservice(lambdaFn string, apiRoute string, bodyJSON map[string]interface{}, idToken string) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    if _, found := bodyJSON["id"]; found {
        delete(bodyJSON, "id")
    }

    client := lambda.NewFromConfig(cfg)

    body, err := json.Marshal(bodyJSON)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    event := map[string]interface{}{
        "httpMethod": "POST",
        "path": fmt.Sprintf("/%v", apiRoute),
        "headers": map[string]string{
            "Content-Type": "application/json",
            "X-IDTOKEN": idToken,
        },
        "body": string(body),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    var response models.Response
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    return response.StatusCode, jsonObject
}

func DeleteFromMicroserviceById(lambdaFn string, apiRoute string, id string, idToken string) (int, map[string]interface{}) {
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    client := lambda.NewFromConfig(cfg)

    event := map[string]interface{}{
        "httpMethod": "DELETE",
        "headers": map[string]string{
            "Content-Type": "application/json",
            "X-IDTOKEN": idToken,
        },
        "path": fmt.Sprintf("/%v/%v", apiRoute, id),
    }

    eventJSON, err := json.Marshal(event)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }
    
    res, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
        FunctionName: aws.String(lambdaFn),
        Payload: eventJSON,
    })

    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    var response models.Response
    
    reader := bytes.NewReader(res.Payload)
    decode := json.NewDecoder(reader)
    err = decode.Decode(&response)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    if response.Body == "404 page not found"{
        return http.StatusNotFound, map[string]interface{}{"data": response.Body}
    }

    var jsonObject map[string]interface{}
    err = json.Unmarshal([]byte(response.Body), &jsonObject)
    if err != nil {
        return http.StatusInternalServerError, map[string]interface{}{"data": err.Error()}
    }

    return response.StatusCode, jsonObject
}
