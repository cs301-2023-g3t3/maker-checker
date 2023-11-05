package main

import (
	"makerchecker-api/configs"
	"makerchecker-api/routes"
)

func init() {
	configs.ConnectToRedis()
	configs.EnvMongoUri()
}

func main() {
    routes.InitRoutes()
}
