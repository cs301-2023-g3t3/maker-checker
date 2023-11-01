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

// Get both maker and checker requests using UserId
func (t MakercheckerController) GetRequestsByUserId(c *gin.Context) {
    userId := c.Param("userId")
    if userId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "userId cannot be empty",
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    makerData := make(map[string]interface{})
    checkerData := make(map[string]interface{})

    makerFilter := bson.M{"makerId": userId}
    makerCursor, err := collection.Find(ctx, makerFilter)
    if err != nil {
        panic(err)
    }

    checkerFilter := bson.M{"checkerId": userId}
    checkerCursor, err := collection.Find(ctx, checkerFilter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: err.Error(),
        })
        return
    }

    defer makerCursor.Close(ctx)
    defer checkerCursor.Close(ctx)

    // Iterate through makerCursor and populate makerData
    for makerCursor.Next(ctx) {
        var makerResult bson.M
        if err := makerCursor.Decode(&makerResult); err != nil {
            c.JSON(http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError,
                Message: err.Error(),
            })
            return
        }
        makerData[generateUniqueKey(makerResult)] = makerResult
    }

    // Iterate through checkerCursor and populate checkerData
    for checkerCursor.Next(ctx) {
        var checkerResult bson.M
        if err := checkerCursor.Decode(&checkerResult); err != nil {
            c.JSON(http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError,
                Message: err.Error(),
            })
            return
        }
        checkerData[generateUniqueKey(checkerResult)] = checkerResult
    }

    if len(makerData) == 0 && len(checkerData) == 0 {
        c.JSON(http.StatusNotFound, models.HttpError{
            Code: http.StatusNotFound,
            Message: "No requests can be found with this userId",
        })
        return
    }

    res := map[string]interface{}{
        "makerRequests": makerData,
        "checkerRequests": checkerData,
    }

    c.JSON(http.StatusOK, res)
}

func (t MakercheckerController) GetPendingApprovalByUserId(c *gin.Context) {
    userId := c.Param("userId")
    if userId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "User Id cannot be empty.",
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"checkerId": userId, "status": "pending"}
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
