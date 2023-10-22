package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

)

type HealthController struct{}

func (h HealthController) CheckHealth(c *gin.Context) {
	log.Println("Checking Health")
	c.String(http.StatusOK, "Success")
}
