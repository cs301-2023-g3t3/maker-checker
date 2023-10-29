package routes

import (
	"context"
	"makerchecker-api/controllers"
	"makerchecker-api/controllers/makerchecker"
	"makerchecker-api/middleware"
	"os"

    "github.com/gin-contrib/cors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func InitRoutes() {
    PORT := os.Getenv("PORT")

    health := new(controllers.HealthController)
    makerchecker := new(makerchecker.MakercheckerController)
    
    router := gin.Default()

    config := cors.DefaultConfig()
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

    v1 := router.Group("/makerchecker")

    healthGroup := v1.Group("/health")
    healthGroup.GET("", health.CheckHealth)

    makercheckerGroup := v1.Group("/record")
    makercheckerGroup.Use(middleware.DecodeJWT())
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)
    makercheckerGroup.GET("/:makercheckerId", makerchecker.GetMakercheckerById)

    makercheckerGroup.GET("/checker/:checkerId", makerchecker.GetByCheckerId)
    makercheckerGroup.GET("/checker/:checkerId/:status", makerchecker.GetByCheckerId)

    makercheckerGroup.GET("/maker/:makerId", makerchecker.GetByMakerId)
    makercheckerGroup.GET("/maker/:makerId/:status", makerchecker.GetByMakerId)

    makercheckerGroup.POST("", makerchecker.CreateMakerchecker)

    makercheckerGroup.PUT("/:makercheckerId/:status", makerchecker.UpdateMakerchecker)

    env := os.Getenv("ENV")
    if env == "lambda" {
        ginLambda = ginadapter.New(router)
        lambda.Start(Handler)
    } else {
        router.Run(":"+ PORT)
    }
}
