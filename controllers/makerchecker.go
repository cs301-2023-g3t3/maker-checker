package controllers

import (
	"context"
	"makerchecker/configs"
	"makerchecker/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MakercheckerController struct{}

var collection *mongo.Collection = configs.OpenCollection(configs.Client, "makerchecker")
var validate = validator.New()

func (t MakercheckerController) GetAllMakercheckers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    options := options.Find()
    cursor, err := collection.Find(ctx, bson.M{}, options)
    if err != nil {
        panic(err)
    }
    
    defer cursor.Close(ctx)

    var makercheckers [] models.Makerchecker
    err = cursor.All(ctx, &makercheckers)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HTTPError{Code: http.StatusInternalServerError, Message: "Failed to retrieve makerchecker requests"})
        return
    }

    c.JSON(http.StatusOK, makercheckers)
}

func (t MakercheckerController) PostMakerchecker (c *gin.Context) {
    data := new(models.Makerchecker)
    err := c.BindJSON(data)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Invalid Makerchecker object"})
        return
    }

    // Validate the Makerchecker object
    if validationErr := validate.Struct(data); validationErr != nil {
        if bodyErr := (data.Data.User == nil && data.Data.Point == nil); !bodyErr {
            c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Invalid Makerchecker object"})
            return
        }
    }
    
    // Validate if the body have either User/Point object
    if bodyErr := (data.Data.User == nil && data.Data.Point == nil); bodyErr {
        c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: "Provide User or Point object"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := collection.InsertOne(ctx, data)

    if err != nil {
        msg := "Failed to insert makerchecker" + err.Error()
        if mongo.IsDuplicateKeyError(err) {
            msg = "Id already exists"
        }
        c.JSON(http.StatusBadRequest, models.HTTPError{Code: http.StatusBadRequest, Message: msg})
        return
    }

    c.JSON(http.StatusCreated, result)
}
