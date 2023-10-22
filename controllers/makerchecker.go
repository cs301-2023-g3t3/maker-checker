package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"makerchecker/configs"
	"makerchecker/middleware"
	"makerchecker/models"
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
            http.StatusInternalServerError, models.HttpResponse{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve makerchecker requests",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(
        http.StatusOK, models.HttpResponse{
            Code: http.StatusOK, 
            Message: "Success",
            Data: map[string]interface{}{"data": makercheckers},
    })
}

func (t MakercheckerController) GetPendingWithCheckerId(c *gin.Context) {
    checkerId := c.Param("checkerId");
    if checkerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpResponse{
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
        c.JSON(http.StatusInternalServerError, models.HttpResponse{
            Code: http.StatusInternalServerError,
            Message: "Failed to retrieve makerchecker requests.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if len(makercheckers) == 0 {
        c.JSON(http.StatusNotFound, models.HttpResponse{
            Code: http.StatusNotFound,
            Message: "Not found",
            Data: map[string]interface{}{"data": "CheckerID: " + checkerId + " not found."},
        })
        return
    }
    
    c.JSON(http.StatusOK, models.HttpResponse{
        Code: http.StatusOK,
        Message: "Success",
        Data: map[string]interface{}{"data": makercheckers},
    })
}

func (t MakercheckerController) GetPendingWithMakerId(c *gin.Context) {
    makerId := c.Param("makerId");
    if makerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpResponse{
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
        c.JSON(http.StatusInternalServerError, models.HttpResponse{
            Code: http.StatusInternalServerError,
            Message: "Failed to retrieve makerchecker requests.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if len(makercheckers) == 0 {
        c.JSON(http.StatusNotFound, models.HttpResponse{
            Code: http.StatusNotFound,
            Message: "Not found",
            Data: map[string]interface{}{"data": "MakerID: " + makerId + " not found."},
        })
        return
    }

    c.JSON(http.StatusOK, models.HttpResponse{
        Code: http.StatusOK,
        Message: "Success",
        Data: map[string]interface{}{"data": makercheckers},
    })
}

func (t MakercheckerController) PostMakerchecker (c *gin.Context) {
    makerchecker := new(models.Makerchecker)
    err := c.BindJSON(makerchecker)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.HttpResponse{
            Code: http.StatusBadRequest,
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    // Validate the Makerchecker object
    if err := validate.Struct(makerchecker); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpResponse{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    // Validate if data is empty
    if len(makerchecker.Data) == 0 {
        c.JSON(http.StatusBadRequest, models.HttpResponse{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": "Data cannot be empty"},
        })
        return
    }

    // Validate if data has an ID
    if _, found := makerchecker.Data["id"]; !found {
        c.JSON(http.StatusBadRequest, models.HttpResponse{
            Code: http.StatusBadRequest, 
            Message: "Invalid Makerchecker object.",
            Data: map[string]interface{}{"data": "Data ID cannot be empty"},
        })
        return
    }
    
    var response models.Response

    // Fetch data from relevant data from microservices
    dataId := fmt.Sprintf("%v", makerchecker.Data["id"])
    dbData := middleware.GetPointById(dataId)
    json.Unmarshal(dbData, &response)

    if response.StatusCode == 404 {
        c.JSON(http.StatusNotFound, models.HttpResponse{
            Code: http.StatusNotFound,
            Message: "Data that you are trying to update does not exist",
            Data: map[string]interface{}{"data": response.Body},
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // add default Status: pending
    makerchecker.Status = "pending"
    makercheckerId := primitive.NewObjectID().Hex()
    makerchecker.MakercheckerId = makercheckerId
    result, err := collection.InsertOne(ctx, makerchecker)

    if err != nil {
        msg := "Failed to insert makerchecker."
        if mongo.IsDuplicateKeyError(err) {
            msg = "MakercheckId already exists."
        }
        c.JSON(http.StatusBadRequest, models.HttpResponse{
            Code: http.StatusBadRequest,
            Message: msg,
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusCreated, models.HttpResponse{
        Code: http.StatusCreated,
        Message: "Success",
        Data: map[string]interface{}{"data": makerchecker, "id": result},
    })
}
