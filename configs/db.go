package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
    mongouri := EnvMongoUri()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.NewClient(options.Client().ApplyURI(mongouri))
    if err != nil {
        log.Fatal("Failed to create client")
    }

    fmt.Printf("Connecting to MongoDB...")
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal("Fail to connect to DB")
    }

    fmt.Println("Success!")

    return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    var collection *mongo.Collection = client.Database("cs301").Collection(collectionName)
    return collection
}
