package routes

import (
	"context"
	"makerchecker-api/controllers"
	"makerchecker-api/controllers/makerchecker"
	permission "makerchecker-api/controllers/permissions"
	// "makerchecker-api/middleware"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	// "github.com/gin-contrib/cors"
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
    permission := new(permission.PermissionController)
    
    router := gin.Default()

 //    config := cors.DefaultConfjjjjjjjjjjig()
	// config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	// config.AllowAllOrigins = true
	// router.Use(cors.New(config))

    v1 := router.Group("/makerchecker")

    healthGroup := v1.Group("/health")
    healthGroup.GET("", health.CheckHealth)

    makercheckerGroup := v1.Group("/record")
    // makercheckerGroup.Use(middleware.DecodeJWT())
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)
    makercheckerGroup.GET("/:makercheckerId", makerchecker.GetMakercheckerById)
    makercheckerGroup.GET("/user/:userId", makerchecker.GetRequestsByUserId)
    makercheckerGroup.GET("/checker/:userId", makerchecker.GetByCheckerId)
    makercheckerGroup.GET("/checker/:userId/:status", makerchecker.GetByCheckerId)
    makercheckerGroup.GET("/maker/:userId", makerchecker.GetByMakerId)
    makercheckerGroup.GET("/maker/:userId/:status", makerchecker.GetByMakerId)

    makercheckerGroup.POST("", makerchecker.CreateMakerchecker)
    makercheckerGroup.POST("/check", makerchecker.CheckMakerchecker)
    makercheckerGroup.PUT("/:id/:status", makerchecker.UpdateMakerchecker)

    permissionGroup := v1.Group("/permission") 
    // permissionGroup.Use(middleware.DecodeJWT())
    permissionGroup.GET("", permission.GetAllPermission)
    permissionGroup.GET("/:id", permission.GetPermissionById)
    permissionGroup.GET("/byRoute", permission.GetPermissionByRoute)
    permissionGroup.POST("", permission.CreateMakercheckerPermission)
    permissionGroup.PUT("/:id", permission.UpdatePermissionById)
    permissionGroup.DELETE("/:id", permission.DeletePermissionById)

    env := os.Getenv("ENV")
    if env == "lambda" {
        ginLambda = ginadapter.New(router)
        lambda.Start(Handler)
    } else {
        router.Run(":"+ PORT)
    }
}
