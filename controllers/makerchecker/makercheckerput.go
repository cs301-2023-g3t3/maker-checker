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
        fmt.Println(responseBody)
        break
    default:
        return http.StatusBadRequest, map[string]interface{}{"data": "Action must be either 'POST', 'PUT', or 'DELETE' only"}
    }

    return statusCode, responseBody
}

func (t MakercheckerController) UpdateMakerchecker (c *gin.Context) {
    id := c.Param("id")
    status := c.Param("status")
    if id == "" || status == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Id and status cannot be empty",
            Data: nil,
        })
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var makerchecker models.Makerchecker

    filter := bson.M{"_id": id}
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

    if status == "cancelled"{
        makerchecker.Status = "cancelled"
    } else if status == "approved" {
        statusCode, responseBody = RequestApproved(lambdaFn, apiRoute, makerchecker.Data, endpointParts[2])
        if statusCode != 200 && statusCode != 201 {
            msg := fmt.Sprint(responseBody)
            if statusCode == 0 {
                statusCode = 500
                msg = "Error retrieving data from the microservices."
            }

            c.JSON(statusCode, models.HttpError{
                Code: statusCode,
                Message: "Error",
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
