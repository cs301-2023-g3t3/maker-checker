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

func (t PermissionController) UpdatePermissionById(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "id cannot be empty",
            Data: nil,
        })
    }

    type UpdatePermission struct {
        Route           string                      `json:"route" bson:"route" validation:"required"`
        Maker          []string                     `json:"maker" bson:"maker" validation:"required"`
        Checker          []string                   `json:"checker" bson:"checker" validation:"required"`
    }

    var permission UpdatePermission
    if err := c.BindJSON(&permission); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid permission object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if err := validate.Struct(permission); err != nil {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid permission object.",
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var updatedData models.Permission
    filter := bson.M{"_id": id}
    err := collection.FindOneAndReplace(ctx, filter, permission).Decode(&updatedData)
    if err != nil {
        msg := "Failed to update permission record"
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

    c.JSON(200, map[string]interface{}{"Updated Makerchecker": permission})
} 
