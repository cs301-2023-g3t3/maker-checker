package configs

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func InitEnvironment() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println(fmt.Sprint(err))
        log.Fatal("Error loading .env file")
    }
}
