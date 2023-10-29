package makerchecker

import (
	"context"
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

func (t MakercheckerController) CheckMakerchecker (c *gin.Context) {
    type ValidMakerchecker struct {
        MakerRole       string                      `json:"makerRole" bson:"makerRole" validate:"required"`
        CheckerRole     string                      `json:"checkerRole" bson:"checkerRole" validate:"required"`
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
    // Validate if a route for makerchecker exists
    permission, found, err := permission.FindPermissionByRoute(requestRoute)
    if !found {
        c.JSON(http.StatusNotFound, models.HttpError{
            Code: http.StatusNotFound,
            Message: "Endpoint route does not allow makerchecker",
            Data: map[string]interface{}{"data": err},
        })
        return
    }

    if validMakerRole := slices.Contains(permission.Maker, validMakerchecker.MakerRole); !validMakerRole {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden, 
            Message: "Maker does not have enough permission to do makerchecker",
            Data: map[string]interface{}{"data": "Invalid maker role"},
        })
        return
    }

    if validCheckerRole := slices.Contains(permission.Checker, validMakerchecker.CheckerRole); !validCheckerRole {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden, 
            Message: "Checker does not have enough permission to do makerchecker",
            Data: map[string]interface{}{"data": "Invalid checker role"},
        })
        return
    }

    statusCode, responseBody := middleware.GetListofUsersWithRoles(permission.Checker) 

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
    permission, found, err := permission.FindPermissionByRoute(requestRoute)
    if !found {
        c.JSON(http.StatusNotFound, models.HttpError{
            Code: http.StatusNotFound,
            Message: "Endpoint route does not allow makerchecker",
            Data: map[string]interface{}{"data": err},
        })
        return
    }

    if validMakerRole := slices.Contains(permission.Maker, makerchecker.MakerRole); !validMakerRole {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden, 
            Message: "Maker does not have enough permission to do makerchecker",
            Data: map[string]interface{}{"data": "Invalid maker role"},
        })
        return
    }

    if validCheckerRole := slices.Contains(permission.Checker, makerchecker.CheckerRole); !validCheckerRole {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden, 
            Message: "Checker does not have enough permission to do makerchecker",
            Data: map[string]interface{}{"data": "Invalid checker role"},
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
    // TODO: check key-value pairs with relevant microservices, i.e. models
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
            Data: nil,
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
