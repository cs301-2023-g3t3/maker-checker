package makerchecker

import (
	"context"
    "fmt"
	"makerchecker-api/middleware"
	"makerchecker-api/models"
	"makerchecker-api/utils"
    "strings"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func RequestApproved(lambdaFn string, apiRoute string, data map[string]interface{}, action string) (int, map[string]interface{}) {
    var statusCode int
    var responseBody map[string]interface{}

    switch action {
    case "PUT":
        statusCode, responseBody = middleware.UpdateWithMicroservice(lambdaFn, apiRoute, data)
        break
    case "POST":
        statusCode, responseBody = middleware.CreateWithMicroservice(lambdaFn, apiRoute, data)
        break
    case "DELETE":
        statusCode, responseBody = middleware.DeleteFromMicroserviceById(lambdaFn, apiRoute, fmt.Sprint(data["id"]))
        break
    default:
        return http.StatusBadRequest, map[string]interface{}{"data": "Action must be either 'POST', 'PUT', or 'DELETE' only"}
    }

    return statusCode, responseBody
}

func (t MakercheckerController) UpdateMakerchecker (c *gin.Context) {
    type RequestBody struct {
        Id      string      `json:"id" validate:"required"`
        Status  string      `json:"status" validate:"required"`
    }

    var requestBody RequestBody
    if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    if err := validate.Struct(requestBody); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    checkerDetails := GetUserDetails(c)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var makerchecker models.Makerchecker

    filter := bson.M{"_id": requestBody.Id}
    err := collection.FindOne(ctx, filter).Decode(&makerchecker)
    if err != nil {
        msg := "Failed to retrieve makerchecker record"
        if err != mongo.ErrNoDocuments {
            msg = "No makerchecker record with makercheckerId"
        }
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: msg,
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if makerchecker.CheckerId != checkerDetails.Id {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden,
            Message: "User is not authorize to approve the request.",
        })
        return
    }

    if makerchecker.Status == "cancelled" || makerchecker.Status == "approved"{
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Bad Request",
            Data: map[string]interface{}{"data": "Request cannot be executed"},
        })
        return
    }

    endpointParts := strings.Split(makerchecker.Endpoint, "/")
    lambdaFn, apiRoute := utils.ProcessMicroserviceTypes(endpointParts[3])

    var statusCode int
    var responseBody map[string]interface{}

    if requestBody.Status == "cancelled"{
        makerchecker.Status = "cancelled"
    } else if requestBody.Status == "approved" {
        statusCode, responseBody = RequestApproved(lambdaFn, apiRoute, makerchecker.Data, endpointParts[2])
        if statusCode != 200 && statusCode != 201 {
            msg := fmt.Sprint(responseBody)
            if statusCode == 0 {
                statusCode = 500
                msg = "Error making request to microservices"
            }

            c.JSON(statusCode, models.HttpError{
                Code: statusCode,
                Message: "Error making request to microservices",
                Data: map[string]interface{}{"data": msg},
            })
            return
        }
        makerchecker.Status = "approved"
    } else {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid type of status",
            Data: map[string]interface{}{"data": "Status must be either 'approved' or 'cancelled'"},
        })
        return
    }

    update := bson.M{"$set": makerchecker}
    _ , err = collection.UpdateOne(ctx, filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Unable to update makerchecker record",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(200, map[string]interface{}{"Updated Makerchecker": makerchecker, "Updated Data": responseBody})
} 
