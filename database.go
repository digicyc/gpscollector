package main

import (
	"context"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)


func MongoClose(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
    defer cancel()

    defer func() {
        if err := client.Disconnect(ctx); err != nil {
            panic(err)
        }
    }()
}

func MongoConnect(uri string)(*mongo.Client, context.Context, context.CancelFunc, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return client, ctx, cancel, err
}


func MongoPing(client *mongo.Client, ctx context.Context) error {
    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        return err
    }
    return nil
}

func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{})(*mongo.InsertOneResult, error) {
    collection := client.Database(dataBase).Collection(col)

    result, err := collection.InsertOne(ctx, doc)

    return result, err
}

func InsertMany(client *mongo.Client, ctx context.Context, dataBase, col string, docs []interface{})(*mongo.InsertManyResult, error) {
    collection := client.Database(dataBase).Collection(col)

    result, err := collection.InsertMany(ctx, docs)

    return result, err
}
