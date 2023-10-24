package makerchecker

import (
	"context"
	"makerchecker-api/middleware"
	"makerchecker-api/models"
	"makerchecker-api/utils"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t MakercheckerController) UpdateMakerchecker (c *gin.Context) {
    makercheckerId := c.Param("makercheckerId")
    status := c.Param("status")
    if makercheckerId == "" || status == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "MakercheckerID and status cannot be empty",
            Data: nil,
        })
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var makerchecker models.Makerchecker

    filter := bson.M{"makercheckerId": makercheckerId}
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

    lambdaFn, apiRoute := utils.ProcessMicroserviceTypes(makerchecker)

    var statusCode int
    var responseBody map[string]interface{}

    switch status {
    case "cancelled":
        makerchecker.Status = "cancelled"
        break
    case "approved":
        if makerchecker.Action == "UPDATE" {
            statusCode, responseBody = middleware.UpdateMicroserviceById(lambdaFn, apiRoute, makerchecker.Data)
        } else {
            statusCode, responseBody = middleware.CreateServiceWithMicroservice(lambdaFn, apiRoute, makerchecker.Data)
        }

        makerchecker.Status = "approved"
        break
    default:
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid type of status",
            Data: map[string]interface{}{"data": "Status must be either 'approved' or 'cancelled'"},
        })
        return
    }

    if statusCode != 200 && statusCode != 201 {
        msg := responseBody["message"]
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
