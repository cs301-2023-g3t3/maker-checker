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

func generateUniqueKey (data bson.M) string {
    if data["makerId"] != 0 {
        return data["makerId"].(string)
    } else {
        return data["checkerId"].(string)
    }
}

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
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Id cannot be empty",
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var makerchecker models.Makerchecker
    filter := bson.M{"_id": id}
    err := collection.FindOne(ctx, filter).Decode(&makerchecker)
    if err != nil {
        msg := "Failed to retrieve makerchecker record"
        if err != mongo.ErrNoDocuments {
            msg = "No makerchecker record with id"
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

func (t MakercheckerController) GetPendingApprovalByMakerId(c *gin.Context) {
    makerId := c.Param("userId")
    if makerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "User Id cannot be empty.",
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"makerId": makerId, "status": "pending"}
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Unable to retrieve data",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    defer cursor.Close(ctx)
    var makerchecker []models.Makerchecker
    err = cursor.All(ctx, &makerchecker)
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve makerchecker requests",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if len(makerchecker) == 0 {
        c.JSON(
            http.StatusNotFound, models.HttpError{
                Code: http.StatusNotFound, 
                Message: "No pending requests found",
        })
        return
    }

    c.JSON(http.StatusOK, makerchecker)
}

func (t MakercheckerController) GetPendingApprovalByCheckerId(c *gin.Context) {
    checkerId := c.Param("userId")
    if checkerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "User Id cannot be empty.",
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"checkerId": checkerId, "status": "pending"}
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Unable to retrieve data",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    defer cursor.Close(ctx)
    var makerchecker []models.Makerchecker
    err = cursor.All(ctx, &makerchecker)
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve makerchecker requests",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if len(makerchecker) == 0 {
        c.JSON(
            http.StatusNotFound, models.HttpError{
                Code: http.StatusNotFound, 
                Message: "No pending requests found",
        })
        return
    }

    c.JSON(http.StatusOK, makerchecker)
}
