package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnect(uri string, base string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	db := client.Database(base)
	if db == nil {
		panic("No collections was found")
	}
	return db
}

// uri:	"mongodb://localhost:27017/"
// db: 	"sacred_base"
// collection: "users"
