package makerchecker

import (
	"context"
	"makerchecker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (t MakercheckerController) GetMakercheckerById(c *gin.Context) {
    makercheckerId := c.Param("makercheckerId")

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

    c.JSON(http.StatusOK, makerchecker)
}

func (t MakercheckerController) GetByCheckerId(c *gin.Context) {
    checkerId := c.Param("checkerId");
    status := c.Param("status")

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

    var filter bson.M
    if status != "" {
        filter = bson.M{"checkerId": checkerId, "status": status}
    } else {
        filter = bson.M{"checkerId": checkerId}
    }

    cursor, err := collection.Find(ctx, filter)
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
            Data: map[string]interface{}{"data": "Records with CheckerID: " + checkerId + " not found."},
        })
        return
    }
    
    c.JSON(http.StatusOK, makercheckers)
}

func (t MakercheckerController) GetByMakerId(c *gin.Context) {
    makerId := c.Param("makerId");
    status := c.Param("status")
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

    var filter bson.M
    if status != "" {
        filter = bson.M{"makerId": makerId, "status": status}
    } else {
        filter = bson.M{"makerId": makerId}
    }

    cursor, err := collection.Find(ctx, filter)
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
            Data: map[string]interface{}{"data": "Records with MakerID: " + makerId + " not found."},
        })
        return
    }

    c.JSON(http.StatusOK, makercheckers)
}
