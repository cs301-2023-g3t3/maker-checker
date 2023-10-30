package makerchecker

import (
	"errors"
	"net/http"
	"slices"
	"makerchecker-api/controllers/permissions"
	"makerchecker-api/configs"
	"makerchecker-api/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type MakercheckerController struct{}


var collection *mongo.Collection = configs.OpenCollection(configs.Client, "makerchecker")
var validate = validator.New()

func Validate(makerRole string, checkerRole string, endpoint string) (int, string, error, *models.Permission) {
    permission, found, err := permission.FindPermissionByEndpoint(endpoint)
    if !found {
        msg := "Endpoint route does not allow makerchecker"
        return http.StatusNotFound, msg, err, nil
    }

    if err != nil {
        return http.StatusInternalServerError, err.Error(), err, nil
    }

    if validMakerRole := slices.Contains(permission.Maker, makerRole); makerRole != "" && !validMakerRole {
        msg := "Maker does not have enough permission to do makerchecker"
        return http.StatusForbidden, msg, errors.New("Invalid maker role"), nil
    }

    if validCheckerRole := slices.Contains(permission.Checker, checkerRole); checkerRole != "" && !validCheckerRole{
        msg := "Checker does not have enough permission to do makerchecker"
        return http.StatusForbidden, msg, errors.New("Checker maker role"), nil
    }

    return http.StatusOK, "Success", nil, &permission
}

