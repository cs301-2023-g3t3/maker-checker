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

//  @Summary        Get all Makerchecker Record
//  @Description    Retrieves a list of Makerchecker records
//  @Tags           makerchecker
//  @Produce        json
//  @Success        200     {array}     models.Makerchecker
//  @Failure        500     {object}    models.HttpError
//  @Router         /record   [get]
func (t MakercheckerController) GetAllMakercheckers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(
            http.StatusInternalServerError, models.HttpError{
                Code: http.StatusInternalServerError, 
                Message: "Failed to retrieve makerchecker requests",
                Data: map[string]interface{}{"data": err.Error()},
        })
        return
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

//  @Summary        Get Makerchecker recordby Id
//  @Description    Retrieve a Makerchecker By ID
//  @Tags           makerchecker
//  @Produce        json
//  @Param          id      path    string  true    "id"
//  @Success        200     {object}    models.Makerchecker
//  @Failure        400     {object}    models.HttpError    "Id cannot be empty"
//  @Failure        404     {object}    models.HttpError    "Record not found with Id"
//  @Failure        500     {object}    models.HttpError
//  @Router         /record/{id}   [get]
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

//  @Summary        Get all Pending Approval as a Maker using MakerID
//  @Description    Retrieves a list of pending approval records
//  @Tags           makerchecker
//  @Produce        json
//  @Success        200     {array}     models.Makerchecker
//  @Failure        400     {object}    models.HttpError    "Maker Id cannot be found in the header provided"
//  @Failure        404     {object}    models.HttpError    "No pending requests found"
//  @Failure        500     {object}    models.HttpError
//  @Router         /record/pending-approve   [get]
func (t MakercheckerController) GetPendingApprovalByMakerId(c *gin.Context) {
    makerId := GetUserDetails(c).Id
    if makerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Maker Id cannot be found in the header provided",
        })
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

//  @Summary        Get all Pending Approval as a Checker using CheckerID
//  @Description    Retrieves a list of pending approval records
//  @Tags           makerchecker
//  @Produce        json
//  @Success        200     {array}     models.Makerchecker
//  @Failure        400     {object}    models.HttpError    "Checker Id cannot be found in the header provided"
//  @Failure        404     {object}    models.HttpError    "No pending requests found"
//  @Failure        500     {object}    models.HttpError
//  @Router         /record/to-approve   [get]
func (t MakercheckerController) GetPendingApprovalByCheckerId(c *gin.Context) {
    checkerId := GetUserDetails(c).Id
    if checkerId == "" {
        c.JSON(http.StatusBadRequest, models.HttpError{
            Code: http.StatusBadRequest,
            Message: "Checker Id cannot be found in the header provided",
        })
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
