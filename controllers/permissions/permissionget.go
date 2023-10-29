package permission

import (
	"context"
	"makerchecker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t PermissionController) GetAllPermission(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve permissions",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }
    
    defer cursor.Close(ctx)

    var permission [] models.Permission
    err = cursor.All(ctx, &permission)
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve permissions",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusOK, permission)
}

func (t PermissionController) GetPermissionById(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(
            http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest, 
                Message: "id cannot be empty",
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var permission models.Permission

    filter := bson.M{"_id": id}
    err := collection.FindOne(ctx, filter).Decode(&permission)
    if err != nil {
        msg := "Failed to retrieve permission"
        statusCode := http.StatusInternalServerError
        if err == mongo.ErrNoDocuments {
            msg = "No permission found with id"
            statusCode = http.StatusNotFound
        }
        c.JSON(statusCode, models.HttpError{
                Code: statusCode, 
                Message: msg,
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }
    
    c.JSON(http.StatusOK, permission)
}
