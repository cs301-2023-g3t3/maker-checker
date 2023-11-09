package main

import (
	"makerchecker-api/configs"
	"makerchecker-api/routes"

	log "github.com/sirupsen/logrus"
)

func init() {
	configs.ConnectToRedis()
	configs.EnvMongoUri()
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
    routes.InitRoutes()
}
