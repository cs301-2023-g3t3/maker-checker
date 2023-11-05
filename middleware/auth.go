package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"makerchecker-api/models"
	"makerchecker-api/utils"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
)

func VerifyUserInfo() gin.HandlerFunc{
    return func (c *gin.Context) {
        userId := c.Param("userId")
        if userId == "" {
            c.JSON(http.StatusBadRequest, models.HttpError{
                Code: http.StatusBadRequest,
                Message: "Id cannot be empty",
            })
            return
        }

		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(http.StatusForbidden, "No Authorization header provided")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			c.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			c.Abort()
			return
		}

        // initialise new cognitoidentityprovider
        sess, err := session.NewSession()
        svc := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion("ap-southeast-1"))
        
        // Create a GetUserInput object and set the AccessToken field
        input := &cognitoidentityprovider.GetUserInput{
            AccessToken: &tokenString,
        }


        // Use the input to call GetUser
        res, err := svc.GetUser(input)
        if err != nil{
            c.JSON(http.StatusInternalServerError, err)
            c.Abort()
            return
        }

        var email string
        for _, attribute := range res.UserAttributes {
            if *attribute.Name == "email" {
            email = *attribute.Value
            break  // Exit the loop after finding the email attribute
        }
}
        
        lambdaFn, apiRoute := utils.ProcessMicroserviceTypes("users")
        statusCode, resObj := GetFromMicroserviceById(lambdaFn, apiRoute, userId)
        if statusCode != 200 {
            c.JSON(http.StatusUnauthorized, models.HttpError{
                Code: http.StatusUnauthorized,
                Message: "Unable to verify the User ID",
                Data: map[string]interface{}{"data": err},
            })
            c.Abort()
            return
        }

        resEmail, ok := resObj["email"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, models.HttpError{
                Code: http.StatusUnauthorized,
                Message: "Unable to verify the User",
                Data: map[string]interface{}{"data": err},
            })
            c.Abort()
            return
        }

        if resEmail != email {
            c.JSON(http.StatusUnauthorized, models.HttpError{
                Code: http.StatusUnauthorized,
                Message: "Unable to verify the User",
            })
            c.Abort()
            return
        }

        c.Set("userDetails", resObj)
        c.Next()
    }
}

type ParsedUserClaim struct {
    Role        string 	`json:"role"`
    Email       string  `json:"email"`
}

func DecodeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(http.StatusForbidden, "No Authorization header provided")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			c.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			c.Abort()
			return
		}
		var keysJWK = os.Getenv("JWT_SECRET")
		setOfKeys, err := jwk.ParseString(keysJWK)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create JWK: %s", err))
			c.Abort()
			return
		}

        // change this for id_token or access_token
		privKey, success := setOfKeys.Get(1)
		if !success {
			c.String(http.StatusInternalServerError, "Could not find key at given index")
			c.Abort()
			return
		}

		pubkey, err := jwk.PublicKeyOf(privKey)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get public key: %s", err.Error()))
			c.Abort()
			return
		}
  
		verifiedToken, err := jws.Verify([]byte(tokenString), jwa.RS256, pubkey)
		if err != nil {
			c.String(http.StatusForbidden, fmt.Sprintf("Failed to verify token from HTTP request: %s", err.Error()))
			c.Abort()
			return
		}

		if err != nil {
			log.Fatal(err)
		}

		var parsedUser ParsedUserClaim
        err = json.Unmarshal([]byte(verifiedToken), &parsedUser)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal data: %s", err.Error()))
			c.Abort()
			return
		}

        parsedUserBytes, err := json.Marshal(parsedUser)
        if err != nil {
            panic(err)
        }

        var userDetails map[string]interface{}
        err = json.Unmarshal(parsedUserBytes, &userDetails)
        if err != nil {
            panic(err)
        }

        c.Set("userDetails", userDetails)
        c.Next()
	}
}