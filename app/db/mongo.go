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
	clientOptions := options.Client().ApplyURI("mongodb://giangnt:giangntxpower@35.247.129.23:27018/dating_clone?authSource=admin")

	// Kết nối đến MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	clientDB := client.Database("dating_clone")
	datasetCollectionName := clientDB.Collection(string(collection.DatasetCollectionName))

	return &MongoDB{
		DB:                    clientDB,
		DatasetCollectionName: datasetCollectionName,
		UsersCollectionName:   clientDB.Collection("users"),
	}, nil
}
