package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

    fmt.Printf("Connecting to MongoDB...\n")
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal("Fail to connect to DB")
    }

    fmt.Println("Pinging server ...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}

    fmt.Printf("Success!\n")

    fmt.Printf("Creating Indexes\n")

    InitIndexes(client)

    return client
}

func InitIndexes(client *mongo.Client) {
    makercheckerCollection := OpenCollection(client, "makerchecker")
    makercheckerIndexModel := mongo.IndexModel{
        Keys: bson.D{{Key: "makercheckerId", Value: -1}},
        Options: options.Index().SetUnique(true),
    }

    makercheckerIndexCreated, err := makercheckerCollection.Indexes().CreateOne(context.Background(), makercheckerIndexModel)
    if err != nil {
        log.Fatal(err)
    }

    permissionCollection := OpenCollection(client, "permission")
    permissionIndexModel := mongo.IndexModel{
        Keys: bson.D{{Key: "route", Value: -1}},
        Options: options.Index().SetUnique(true),
    }

    permissionIndexCreated, err := permissionCollection.Indexes().CreateOne(context.Background(), permissionIndexModel)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created Makerchecker Index %s\n", makercheckerIndexCreated)
    fmt.Printf("Created Permission Index %s\n", permissionIndexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    var collection *mongo.Collection = client.Database("cs301").Collection(collectionName)
    return collection
}
