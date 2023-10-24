package makerchecker

import (
	"context"
	"fmt"
	"makerchecker-api/middleware"
	"makerchecker-api/models"
	"makerchecker-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t MakercheckerController) CreateMakerchecker (c *gin.Context) {
    makerchecker := new(models.Makerchecker)
    err := c.BindJSON(makerchecker)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    // Validate if data is empty
    if len(makerchecker.Data) == 0 {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": "Data cannot be empty"},
        })
        return
    }

    // Get relevant Lambda Function and API Routes
    lambdaFn, apiRoute := utils.ProcessMicroserviceTypes(*makerchecker)
    
    if lambdaFn == "Error" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": "Database field must be of 'users' or 'points'."},
        })
        return
    }

    var data map[string]interface{}
    
    switch makerchecker.Action {
    case "UPDATE" :
        if _, found := makerchecker.Data["id"]; !found {
            c.JSON(http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest,
                Message: "Data ID cannot be empty",
                Data: nil,
            })
            return
        }
        dataId := fmt.Sprint(makerchecker.Data["id"])

        // Fetch data from relevant data from microservices
        statusCode, responseBody := middleware.GetFromMicroserviceById(lambdaFn, apiRoute, dataId) 

        // Error fetching data
        if statusCode != 200 {
            if statusCode == 0 {
                c.JSON(statusCode, models.HttpError{
                    Code: 500,
                    Message: "Internal Server Error",
                    Data: map[string]interface{}{"data": "Error retrieving data from the microservices. Database is not functioning."},
                })
                return
            }

            c.JSON(statusCode, models.HttpError{
                Code: statusCode,
                Message: "Error",
                Data: map[string]interface{}{"data": responseBody["message"]},
            })
            return
        }

        statusCode, data = utils.GetDifferences(responseBody, makerchecker.Data)

        if statusCode != 200 {
            c.JSON(statusCode, models.HttpError{
                Code: statusCode,
                Message: "Error",
                Data: data,
            })
            return
        }
        makerchecker.Data = data
        break
    case "CREATE":
        if makerchecker.Database != "users" {
            c.JSON(http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest,
                Message: "Invalid action",
                Data: nil,
            })
            return
        }
        break
    default:
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid action",
            Data: nil,
        })
        return
    }


    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    makerchecker.Status = "pending" // add default Status: pending
    makerchecker.MakercheckerId = primitive.NewObjectID().Hex() // add makercheckerId ObjectKey
    result, err := collection.InsertOne(ctx, makerchecker)

    if err != nil {
        msg := "Failed to insert makerchecker."
        if mongo.IsDuplicateKeyError(err) {
            msg = "MakercheckId already exists."
        }
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: msg,
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusCreated, map[string]interface{}{"result": makerchecker, "insertedId": result.InsertedID})
}
