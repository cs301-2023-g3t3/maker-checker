package configs

import (
    "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoUri() string {
	var err error

	env := os.Getenv("ENV")
	if env != "lambda" {
        err = godotenv.Load()
        switch os.Getenv("DB_TYPE") {
            case "local":
                err = godotenv.Load(".env.local")
                return os.Getenv("MONGOURI")
            case "live":
                err = godotenv.Load(".env.live")
            }
            if err != nil {
                log.Fatal("Error loading .env file")
            }
    }
    user := os.Getenv("MONGO_USERNAME")
    pass := os.Getenv("MONGO_PASSWORD")
    host := os.Getenv("MONGO_HOST")

    mongouri := fmt.Sprintf("mongodb+srv://%s:%s@%s.mongodb.net/?retryWrites=true&w=majority", user, pass, host)
    return mongouri
}
