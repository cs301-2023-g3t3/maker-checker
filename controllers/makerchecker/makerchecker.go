package makerchecker

import (
	"errors"
	"makerchecker-api/configs"
	"makerchecker-api/controllers/permissions"
	"makerchecker-api/models"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type MakercheckerController struct{}


var collection *mongo.Collection = configs.OpenCollection(configs.Client, "makerchecker")
var validate = validator.New()

func Validate(makerRole float64, checkerRole float64, endpoint string) (int, string, error, *models.Permission) {
    permission, found, err := permission.FindPermissionByEndpoint(endpoint)
    if !found {
        msg := "Endpoint route does not allow makerchecker"
        return http.StatusNotFound, msg, err, nil
    }

    if err != nil {
        return http.StatusInternalServerError, err.Error(), err, nil
    }

    if validMakerRole := slices.Contains(permission.Maker, makerRole); makerRole != -1 && !validMakerRole {
        msg := "Maker does not have enough permission to do makerchecker"
        return http.StatusForbidden, msg, errors.New("Invalid maker role"), nil
    }

    if validCheckerRole := slices.Contains(permission.Checker, checkerRole); checkerRole != -1 && !validCheckerRole{
        msg := "Checker does not have enough permission to do makerchecker"
        return http.StatusForbidden, msg, errors.New("Checker maker role"), nil
    }

    return http.StatusOK, "Success", nil, &permission
}

type FormattedUserDetails struct {
    Email   string  `json:"email"`
    Id      string  `json:"id"`
    Role    float64 `json:"role"`
}

func GetUserDetails(c *gin.Context) *FormattedUserDetails {
    userDetails, ok := c.Get("userDetails")
    if !ok {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Error",
        })
        c.Abort()
        return nil
    }

    userDetailsObj, ok := userDetails.(map[string]interface{})
    if !ok {
        c.JSON(http.StatusInternalServerError, models.HttpError{
            Code: http.StatusInternalServerError,
            Message: "Error",
        }) 
        c.Abort()
        return nil
    }


    formatUserObj := new(FormattedUserDetails)

    for k, v := range userDetailsObj {
        switch k {
        case "email":
            formatUserObj.Email = v.(string)
            break
        case "user_id":
            formatUserObj.Id = v.(string)
            break
        case "cognito:groups":
            temp, ok := v.([]interface{})
            if !ok {
                c.JSON(http.StatusInternalServerError, models.HttpError{
                    Code: http.StatusInternalServerError,
                    Message: "Unable to retrieve role from Cognito:Groups",
                })
                c.Abort()
                return nil
            }

            val, err := strconv.ParseFloat(temp[0].(string), 64)
            if err != nil {
                c.JSON(http.StatusInternalServerError, models.HttpError{
                    Code: http.StatusInternalServerError,
                    Message: "Error",
                })
                c.Abort()
                return nil
            }
            formatUserObj.Role = val
            break
        }
    }

    return formatUserObj
}
