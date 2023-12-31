package permission

import (
	"context"
	"encoding/json"
	"makerchecker-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindPermissionByEndpoint(endpoint string) (models.Permission, bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var permission models.Permission

    filter := bson.M{"endpoint": endpoint}
    err := collection.FindOne(ctx, filter).Decode(&permission)
    if err != nil {
        return permission, false, err
    }

    return permission, true, nil
}

//  @Summary        Get all Makerchecker Permission
//  @Description    Retrieves a list of makerchecker permission
//  @Tags           permission
//  @Produce        json
//  @Success        200     {array}     models.Permission
//  @Failure        500     {object}    models.HttpError
//  @Router         /permission   [get]
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

//  @Summary        Get Makerchecker permission by Id
//  @Description    Retrieve a Makerchecker permission by Id
//  @Tags           permission
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     {object}    models.Permission
//  @Failure        400     {object}    models.HttpError    "Id cannot be empty"
//  @Failure        404     {object}    models.HttpError    "No permission can be found with Id"
//  @Failure        500     {object}    models.HttpError
//  @Router         /permission/{id}   [get]
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

//  @Summary        Get Makerchecker Permission by Endpoint route
//  @Description    Retrieve a Makerchecker Permission by Endpoint route
//  @Tags           permission
//  @Produce        json
//  @Param          endpoint      body  string  true    "endpoint"
//  @Success        200     {array}     models.Permission
//  @Failure        400     {object}    models.HttpError    "Endpoint cannot be empty"
//  @Failure        404     {object}    models.HttpError    "Makerchecker permission cannot be found with endpoint route"
//  @Failure        500     {object}    models.HttpError
//  @Router         /permission/by-endpoint   [get]
func (t PermissionController) GetPermissionByEndpoint (c *gin.Context) {
    var requestBody map[string]interface{}
    err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
    endpoint, ok := requestBody["endpoint"]
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Error",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
    }

    if !ok {
        c.JSON(
            http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest, 
                Message: "'endpoint' must be in the request body",
        })
        return
    }

    endpointStr, ok := endpoint.(string)
    if !ok {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "'endpoint' must be a string",
        })
        return
    }

    permission, found, err := FindPermissionByEndpoint(endpointStr)
    if !found {
        c.JSON(http.StatusNotFound, models.HttpError{
            Code: http.StatusNotFound,
            Message: "Error",
            Data: map[string]interface{}{"data": "Permission not found"},
        })
        return
    }

    c.JSON(http.StatusOK, permission)
}
