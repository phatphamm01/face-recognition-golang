package db

import (
	"context"
	"face-recognition-golang/db/collection"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB                    *mongo.Database
	DatasetCollectionName *mongo.Collection
	UsersCollectionName   *mongo.Collection
}

var Client *MongoDB

func NewMongoDB() (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI("mongodb://34.143.159.142:27018")

	// Kết nối đến MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	clientDB := client.Database("dating")
	datasetCollectionName := clientDB.Collection(string(collection.DatasetCollectionName))

	return &MongoDB{
		DB:                    clientDB,
		DatasetCollectionName: datasetCollectionName,
		UsersCollectionName:   clientDB.Collection("users"),
	}, nil
}
