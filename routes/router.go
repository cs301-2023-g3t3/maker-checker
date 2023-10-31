package routes

import (
	"context"
	"makerchecker-api/controllers"
	"makerchecker-api/controllers/makerchecker"
	permission "makerchecker-api/controllers/permissions"
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
    
    router := gin.Default()
    router.Use(cors.Default())

    v1 := router.Group("/makerchecker")

    healthGroup := v1.Group("/health")
    healthGroup.GET("", health.CheckHealth)
    
    verifyGroup := v1.Group("/verify")
    verifyGroup.POST("", makerchecker.CheckMakerchecker)

    makercheckerGroup := v1.Group("/record")
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)
    makercheckerGroup.GET("/:id", makerchecker.GetMakercheckerById)
    makercheckerGroup.GET("/user/:userId", makerchecker.GetRequestsByUserId)
    // makercheckerGroup.GET("/checker/:userId", makerchecker.GetByCheckerId)
    // makercheckerGroup.GET("/checker/:userId/:status", makerchecker.GetByCheckerId)
    // makercheckerGroup.GET("/maker/:userId", makerchecker.GetByMakerId)
    // makercheckerGroup.GET("/maker/:userId/:status", makerchecker.GetByMakerId)

    makercheckerGroup.POST("", makerchecker.CreateMakerchecker)
    makercheckerGroup.PUT("/:id/:status", makerchecker.UpdateMakerchecker)

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
