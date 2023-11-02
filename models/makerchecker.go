package models

import (
	"makerchecker-api/configs"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type Makerchecker struct {
    Id              string                      `json:"_id" bson:"_id"`
    MakerId         string                      `json:"makerId" bson:"makerId"`
    MakerEmail      string                      `json:"makerEmail" bson:"makerEmail"`
    CheckerId       string                      `json:"checkerId" bson:"checkerId" validate:"required"`
    CheckerEmail    string                      `json:"checkerEmail" bson:"checkerEmail"`
    Endpoint        string                      `json:"endpoint" bson:"endpoint" validate:"required"`
    Status          string                      `json:"status" bson:"status"`
    Data            map[string]interface{}      `json:"data" bson:"data" validate:"required"`
}


var collection *mongo.Collection = configs.OpenCollection(configs.Client, "makerchecker")
var validate = validator.New()
