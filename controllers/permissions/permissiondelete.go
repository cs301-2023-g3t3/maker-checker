package permission

import (
	"context"
	"makerchecker-api/models"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t PermissionController) DeletePermissionById (c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "id cannot be empty",
        })
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    var permission models.Permission
    filter := bson.M{"_id": id}
    err := collection.FindOneAndDelete(ctx, filter).Decode(&permission)
    if err != nil {
        msg := "Failed to delete permission record"
        statusCode := http.StatusBadRequest
        if err == mongo.ErrNoDocuments {
            msg = "No permission found with id"
            statusCode = http.StatusNotFound
        }
        c.JSON(statusCode, models.HttpError{
            Code: statusCode,
            Message: msg,
            Data: map[string]interface{}{"data": err},
        })
        return
    }

    c.JSON(200, map[string]interface{}{"Deleted Permission": permission})
}
