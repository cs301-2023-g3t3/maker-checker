package controllers

import (
	"context"
	"makerchecker/configs"
	"makerchecker/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MakercheckerController struct{}

var makercheckerCollection *mongo.Collection = configs.OpenCollection(configs.Client, "makerchecker")

func (t MakercheckerController) GetAllMakercheckers(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    options := options.Find()
    cursor, err := makercheckerCollection.Find(ctx, bson.M{}, options)
    if err != nil {
        panic(err)
    }
    
    defer cursor.Close(ctx)

    var makercheckers [] models.Makerchecker
    err = cursor.All(ctx, &makercheckers)
    if err != nil {
        c.JSON(http.StatusInternalServerError, "Failed to retrieve makerchecker requests")
        return
    }

    c.JSON(http.StatusOK, makercheckers)
}
