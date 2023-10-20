package configs

import (
    "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoUri() string {
    err := godotenv.Load()
    if err != nil && os.Getenv("ENV") == "debug" {
        log.Fatal("Error loading .env file")
    }
    
    if os.Getenv("ENV") == "debug" {
        return os.Getenv("MONGOURI")
    } else {
        user := os.Getenv("MONGO_USERNAME")
        pass := os.Getenv("MONGO_PASSWORD")
        host := os.Getenv("MONGO_HOST")
        port := os.Getenv("MONGO_PORT")

        mongouri := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
        return mongouri
    }
}
