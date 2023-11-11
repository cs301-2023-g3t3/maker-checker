package routes

import (
	"context"
	"makerchecker-api/controllers"
	"makerchecker-api/controllers/makerchecker"
	permission "makerchecker-api/controllers/permissions"
	"makerchecker-api/middleware"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
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
    
    // router := gin.Default()
    router := gin.New()
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
    config.AddAllowHeaders("Authorization", "X-IDTOKEN")
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
    router.Use(middleware.LoggingMiddleware())


    v1 := router.Group("/makerchecker")

    healthGroup := v1.Group("/health")
    healthGroup.GET("", health.CheckHealth)
    
    verifyGroup := v1.Group("/verify")
    verifyGroup.Use(middleware.DecodeJWT())
    verifyGroup.POST("", makerchecker.CheckMakerchecker)

    makercheckerGroup := v1.Group("/record")
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)
    makercheckerGroup.Use(middleware.DecodeJWT())
    makercheckerGroup.GET("/:id", makerchecker.GetMakercheckerById)
    makercheckerGroup.GET("/pending-approve", makerchecker.GetPendingApprovalByMakerId)
    makercheckerGroup.GET("/to-approve", makerchecker.GetPendingApprovalByCheckerId)
    makercheckerGroup.POST("", makerchecker.CreateMakerchecker)
    makercheckerGroup.PUT("", makerchecker.UpdateMakerchecker)

    permissionGroup := v1.Group("/permission") 
    permissionGroup.GET("", permission.GetAllPermission)
    permissionGroup.GET("/:id", permission.GetPermissionById)
    permissionGroup.GET("/by-endpoint", permission.GetPermissionByEndpoint)
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
