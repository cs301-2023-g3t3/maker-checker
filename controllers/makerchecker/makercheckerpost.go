package makerchecker

import (
	"context"
	"errors"
	"fmt"
	"makerchecker-api/controllers/permissions"
	"makerchecker-api/middleware"
	"makerchecker-api/models"
	"makerchecker-api/utils"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Validate(makerRole string, checkerRole string, endpoint string) (int, string, error, *models.Permission) {
    permission, found, err := permission.FindPermissionByEndpoint(endpoint)
    if !found {
        msg := "Endpoint route does not allow makerchecker"
        return http.StatusNotFound, msg, err, nil
    }

    if err != nil {
        return http.StatusInternalServerError, err.Error(), err, nil
    }

    if validMakerRole := slices.Contains(permission.Maker, makerRole); !validMakerRole {
        msg := "Maker does not have enough permission to do makerchecker"
        return http.StatusForbidden, msg, errors.New("Invalid maker role"), nil
    }

    if validCheckerRole := slices.Contains(permission.Checker, checkerRole); checkerRole != "" && !validCheckerRole{
        msg := "Checker does not have enough permission to do makerchecker"
        return http.StatusForbidden, msg, errors.New("Checker maker role"), nil
    }

    return http.StatusOK, "Success", nil, &permission
}

func (t MakercheckerController) CheckMakerchecker (c *gin.Context) {
    type ValidMakerchecker struct {
        MakerRole       string                      `json:"makerRole" bson:"makerRole" validate:"required"`
        Endpoint        string                      `json:"endpoint" bson:"endpoint" validate:"required"`
    }

    var validMakerchecker ValidMakerchecker

    if err := c.BindJSON(&validMakerchecker); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Error",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if err := validate.Struct(validMakerchecker); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Invalid data to check for makerchecker",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    requestRoute := validMakerchecker.Endpoint

    statusCode, body, err, permission := Validate(validMakerchecker.MakerRole, "", requestRoute)
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
    makerchecker := new(models.Makerchecker)
    
    if err := c.BindJSON(makerchecker); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if err := validate.Struct(makerchecker); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    requestRoute := makerchecker.Endpoint
    // Validate if a route for makerchecker exists
    statusCode, body, err, _ := Validate(makerchecker.MakerRole, makerchecker.CheckerRole, requestRoute)
    if statusCode != http.StatusOK {
        c.JSON(statusCode, models.HttpError{
            Code: statusCode,
            Message: body,
            Data: map[string]interface{}{"data": err},
        })
        return
    }

    // Get relevant Lambda Function and API Routes
    endpointParts := strings.Split(makerchecker.Endpoint, "/")
    lambdaFn, apiRoute := utils.ProcessMicroserviceTypes(endpointParts[3])
    
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
        if _, found := makerchecker.Data["id"]; !found {
            c.JSON(http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest,
                Message: "Data ID cannot be empty",
                Data: nil,
            })
            return
        }
        dataId := fmt.Sprint(makerchecker.Data["id"])

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

    makerchecker.Status = "pending" // add default Status: pending
    makerchecker.Id = primitive.NewObjectID().Hex() // add makercheckerId ObjectKey

   // TODO: trigger send email

    _, err = collection.InsertOne(ctx, makerchecker)

    if err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message:  "Failed to insert makerchecker.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusCreated, map[string]interface{}{"result": makerchecker})
}
