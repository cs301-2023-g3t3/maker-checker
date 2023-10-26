package main

import (
	"makerchecker-api/configs"
	"makerchecker-api/routes"
)

func init() {
	configs.EnvMongoUri()
}

func main() {
    routes.InitRoutes()
}
