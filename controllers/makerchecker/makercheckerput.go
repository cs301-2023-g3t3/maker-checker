package makerchecker

import (
	// "context"
	// "fmt"
	// "makerchecker-api/middleware"
	"context"
	"fmt"
	"makerchecker-api/models"
	"time"

	// "makerchecker-api/utils"
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t MakercheckerController) UpdateMakerchecker (c *gin.Context) {
    makercheckerId := c.Param("makercheckerId")
    status := c.Param("status")
    if makercheckerId == "" || status == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "MakercheckerID and status cannot be empty",
            Data: nil,
        })
    }

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

    switch status {
    case "cancel":
        makerchecker.Status = "cancel"
    case "ok":
        makerchecker.Status = "ok"
        fmt.Println("Update DB!")    
    default:
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid type of status",
            Data: map[string]interface{}{"data": "Status must be either 'cancel' or 'ok'"},
        })
        return
    }

    if status == "cancel" {
        makerchecker.Status = "cancel"
    } else {
        fmt.Println(status)
    }
    
    update := bson.M{"$set": makerchecker}
    _ , err = collection.UpdateOne(ctx, filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Unable to update makerchecker record",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(200, map[string]interface{}{"Updated Data": makerchecker})
} 
