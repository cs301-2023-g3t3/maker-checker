package routes


import (
	"os"
	"makerchecker-api/controllers"
	"makerchecker-api/controllers/makerchecker"

	"github.com/gin-gonic/gin"
)

func InitRoutes() {
    PORT := os.Getenv("PORT")

    health := new(controllers.HealthController)
    makerchecker := new(makerchecker.MakercheckerController)

    router := gin.Default()
    
    v1 := router.Group("/api/v1")

    healthGroup := v1.Group("/health")
    healthGroup.GET("", health.CheckHealth)

    makercheckerGroup := v1.Group("/makerchecker")
    makercheckerGroup.GET("", makerchecker.GetAllMakercheckers)
    makercheckerGroup.GET("/:makercheckerId", makerchecker.GetMakercheckerById)

    makercheckerGroup.GET("/checker/:checkerId", makerchecker.GetByCheckerId)
    makercheckerGroup.GET("/checker/:checkerId/:status", makerchecker.GetByCheckerId)

    makercheckerGroup.GET("/maker/:makerId", makerchecker.GetByMakerId)
    makercheckerGroup.GET("/maker/:makerId/:status", makerchecker.GetByMakerId)

    makercheckerGroup.POST("", makerchecker.CreateMakerchecker)

    makercheckerGroup.PUT("/:makercheckerId/:status", makerchecker.UpdateMakerchecker)

    router.Run(":"+ PORT)
}
