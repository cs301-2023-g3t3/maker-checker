package routes


import (
	"os"
	"makerchecker-api/controllers/makerchecker"

	"github.com/gin-gonic/gin"
)

func InitRoutes() {
    PORT := os.Getenv("PORT")

    makerchecker := new(controllers.MakercheckerController)

    router := gin.Default()
    
    v1 := router.Group("/api/v1")

    makercheckerGroup := v1.Group("/makerchecker")
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)
    makercheckerGroup.GET("/pending/checker/:checkerId", makerchecker.GetPendingWithCheckerId)
    makercheckerGroup.GET("/pending/maker/:makerId", makerchecker.GetPendingWithMakerId)
    makercheckerGroup.POST("", makerchecker.PostMakerchecker)

    router.Run(":"+ PORT)
}
