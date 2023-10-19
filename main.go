package main

import (
	"makerchecker/configs"
	"makerchecker/routes"
)

func init() {
    configs.InitEnvironment()
}
func main() {
    routes.InitRoutes()
}
