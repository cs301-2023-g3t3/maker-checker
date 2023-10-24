package controllers

import (
	"context"
	"makerchecker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

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
