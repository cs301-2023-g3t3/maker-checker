package makerchecker

import (
	"context"
	"fmt"
	"makerchecker-api/middleware"
	"makerchecker-api/models"
	"makerchecker-api/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t MakercheckerController) CheckMakerchecker (c *gin.Context) {
    userDetails := GetUserDetails(c)

    type RequestBody struct {
        Endpoint    string      `json:"endpoint" validate:"required"`
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

    requestRoute := requestBody.Endpoint

    userRole := userDetails.Role

    statusCode, body, err, permission := Validate(userRole, -1, requestRoute)
    if statusCode != http.StatusOK {
        c.JSON(statusCode, models.HttpError{
            Code: statusCode,
            Message: body,
            Data: map[string]interface{}{"data": err},
        })
        return
    }

    statusCode, responseBody := middleware.GetListofUsersWithRolesWithMicroservice(permission.Checker) 

    if statusCode != 200 {
        msg := fmt.Sprint(responseBody)
        if statusCode == 0 {
            statusCode = 500
            msg = "Error retrieving data from user microservice."
        }

        c.JSON(statusCode, models.HttpError{
            Code: statusCode,
            Message: "Error",
            Data: map[string]interface{}{"data": msg},
        })
        return
    }

    c.JSON(http.StatusOK, responseBody)
}

func (t MakercheckerController) CreateMakerchecker (c *gin.Context) {
    var reqBody models.Makerchecker
    if err := c.BindJSON(&reqBody); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    if err := validate.Struct(reqBody); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: err.Error(),
        })
        return
    }

    makerDetails := GetUserDetails(c)
    makerRole := makerDetails.Role

    // Get checker details
    lambdaFn, apiRoute := utils.ProcessMicroserviceTypes("users")
    statusCode, checkerDetails := middleware.GetFromMicroserviceById(lambdaFn, apiRoute, reqBody.CheckerId)
    if statusCode != 200 {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Checker User ID is not verified.",
        })
        return
    }
    checkerRole, ok := checkerDetails["role"].(float64)
    if !ok {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Something went wrong.",
        })
        return
    }

    requestRoute := reqBody.Endpoint
    // Validate if a route for makerchecker exists
    statusCode, body, err, _ := Validate(makerRole, checkerRole, requestRoute)
    if statusCode != http.StatusOK {
        c.JSON(statusCode, models.HttpError{
            Code: statusCode,
            Message: body,
            Data: map[string]interface{}{"data": err},
        })
        return
    }

    // Get relevant Lambda Function and API Routes
    endpointParts := strings.Split(reqBody.Endpoint, "/")
    lambdaFn, apiRoute = utils.ProcessMicroserviceTypes(endpointParts[3])
    if lambdaFn == "Error" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": "Database field must be of 'users' or 'points'."},
        })
        return
    }

    switch endpointParts[2] {
    case "PUT": case "DELETE":
        if _, found := reqBody.Data["id"]; !found {
            c.JSON(http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest,
                Message: "Data ID cannot be empty",
                Data: nil,
            })
            return
        }
        dataId := fmt.Sprint(reqBody.Data["id"])

        // Check if data exists in relevant database
        statusCode, responseBody := middleware.GetFromMicroserviceById(lambdaFn, apiRoute, dataId) 

        // Error fetching data
        if statusCode != 200 {
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
        break
    case "POST":
        break
    default:
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid action",
            Data: map[string]interface{}{"data": "Action should be 'POST', 'PUT', or 'DELETE' only"},
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    makerEmail := makerDetails.Email

    checkerEmail, ok := checkerDetails["email"].(string)
    res := middleware.TriggerMessageQueueToEmail(makerEmail, checkerEmail)
    if res != "" {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: res,
        })
        return 
    }

    makerId := makerDetails.Id

    reqBody.Id = primitive.NewObjectID().Hex() // add makercheckerId ObjectKey
    reqBody.Status = "pending" // add default Status: pending
    reqBody.MakerId = makerId
    reqBody.MakerEmail = makerEmail
    reqBody.CheckerEmail = checkerEmail

    _, err = collection.InsertOne(ctx, reqBody)

    if err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message:  "Failed to insert makerchecker.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusCreated, map[string]interface{}{"result": reqBody})
}
