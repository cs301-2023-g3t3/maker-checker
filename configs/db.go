package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
    mongouri := os.Getenv("MONGOURI")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    fmt.Println(mongouri)

    client, err := mongo.NewClient(options.Client().ApplyURI(mongouri))
    if err != nil {
        log.Fatal("Failed to create client")
    }

    fmt.Printf("Connecting to MongoDB...")
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal("Fail to connect to DB")
    }

    InitIndexes(client)

    fmt.Println("Success!")

    return client
}

func InitIndexes(client *mongo.Client) {
    makercheckerCollection := OpenCollection(client, "makerchecker")
    makercheckerModel := mongo.IndexModel {
        Keys: bson.D{{Key:"_id", Value: 1}},
        Options: options.Index().SetUnique(true),
    }

    makercheckerIndexCreated, err := makercheckerCollection.Indexes().CreateOne(context.Background(), makercheckerModel)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created Makerchecker Index %s\n", makercheckerIndexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    var collection *mongo.Collection = client.Database("cs301").Collection(collectionName)
    return collection
}
