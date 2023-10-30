package permission

import (
	"context"
	"makerchecker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t PermissionController) CreateMakercheckerPermission(c *gin.Context) {
    var permission models.Permission
    
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

    if len(permission.Maker) == 0 || len(permission.Checker) == 0 {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid permission object.",
            Data: map[string]interface{}{"data": "Maker and checker cannot be empty"},
        })
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    permission.Id = primitive.NewObjectID().Hex()
    _, err := collection.InsertOne(ctx, permission)

    if err != nil {
        msg := "Failed to insert permission record."
        if mongo.IsDuplicateKeyError(err) {
            msg = "Endpoint already exists."
        }
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: msg,
            Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    c.JSON(http.StatusCreated, map[string]interface{}{"result": permission})
}
