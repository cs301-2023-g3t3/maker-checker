package routes


import (
	"os"
	"makerchecker/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes() {
    PORT := os.Getenv("SERVER_PORT")

    makerchecker := new(controllers.MakercheckerController)

    router := gin.Default()
    
    v1 := router.Group("/api/v1")

    makercheckerGroup := v1.Group("/makerchecker")
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)

    router.Run(":"+ PORT)
}
