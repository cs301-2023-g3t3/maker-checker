package controllers

import (
	"context"
	"fmt"
	"makerchecker/configs"
	"makerchecker/middleware"
	"makerchecker/models"
	"makerchecker/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MakercheckerController struct{}


var collection *mongo.Collection = configs.OpenCollection(configs.Client, "makerchecker")
var validate = validator.New()

func (t MakercheckerController) GetAllMakercheckers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        panic(err)
    }
    
    defer cursor.Close(ctx)

    var makercheckers [] models.Makerchecker
    err = cursor.All(ctx, &makercheckers)
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve makerchecker requests",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusOK, makercheckers)
}

func (t MakercheckerController) GetPendingWithCheckerId(c *gin.Context) {
    checkerId := c.Param("checkerId");
    if checkerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Error",
            Data: map[string]interface{}{"data": "CheckerId parameter cannot be empty"},
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"checkerId": checkerId})
    if err != nil {
        panic(err)
    }
    
    defer cursor.Close(ctx)

    var makercheckers [] models.Makerchecker
    err = cursor.All(ctx, &makercheckers)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Failed to retrieve makerchecker requests.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if len(makercheckers) == 0 {
        c.JSON(http.StatusNotFound, models.HttpError{
            Code: http.StatusNotFound,
            Message: "Not found",
            Data: map[string]interface{}{"data": "CheckerID: " + checkerId + " not found."},
        })
        return
    }
    
    c.JSON(http.StatusOK, makercheckers)
}

func (t MakercheckerController) GetPendingWithMakerId(c *gin.Context) {
    makerId := c.Param("makerId");
    if makerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Error",
            Data: map[string]interface{}{"data": "MakerId parameter cannot be empty"},
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{"makerId": makerId})
    if err != nil {
        panic(err)
    }
    
    defer cursor.Close(ctx)

    var makercheckers [] models.Makerchecker
    err = cursor.All(ctx, &makercheckers)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Failed to retrieve makerchecker requests.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if len(makercheckers) == 0 {
        c.JSON(http.StatusNotFound, models.HttpError{
            Code: http.StatusNotFound,
            Message: "Not found",
            Data: map[string]interface{}{"data": "MakerID: " + makerId + " not found."},
        })
        return
    }

    c.JSON(http.StatusOK, makercheckers)
}

func (t MakercheckerController) PostMakerchecker (c *gin.Context) {
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

    // Validate the Makerchecker object
    if err := validate.Struct(makerchecker); err != nil {
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

    // Validate if data has an ID
    if _, found := makerchecker.Data["id"]; !found {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": "Data ID cannot be empty"},
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
    
    // for actions that needs existing data such as UPDATE
    if makerchecker.Action == "UPDATE" {
        dataId := fmt.Sprint(makerchecker.Data["id"])
        statusCode, responseBody := middleware.GetFromMicroserviceById(lambdaFn, apiRoute, dataId) // Fetch data from relevant data from microservices

        if statusCode != 200 {
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
