package permission

import (
	"makerchecker-api/configs"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type PermissionController struct{}


var collection *mongo.Collection = configs.OpenCollection(configs.Client, "permission")
var validate = validator.New()
