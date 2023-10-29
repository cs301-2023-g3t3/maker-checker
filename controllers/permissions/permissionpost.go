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
    userDetails, ok := c.Get("userDetails")
    if !ok {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden,
            Message: "Unable to retrieve user details",
            Data: map[string]interface{}{"data": ok},
        })
        return
    }

    userDetailsMap, ok := userDetails.(map[string]interface{})
    if !ok {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Failed to convert user details to map",
            Data: map[string]interface{}{"error": "Type assertion failed"},
        })
        return
    }

    role, roleExists := userDetailsMap["role"]
    if !roleExists {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden,
            Message: "Role information not found in user details",
        })
        return
    }

    if role != "Owner" {
        c.JSON(http.StatusForbidden, models.HttpError{
            Code: http.StatusForbidden,
            Message: "Not enough permission to create permission record",
        })
        return
    }

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
            msg = "Route already exists."
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
