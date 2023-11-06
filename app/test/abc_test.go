package socketRouter

import (
	"context"
	"face-recognition-golang/db"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test(t *testing.T) {
	objID, _ := primitive.ObjectIDFromHex("65039bc7c0691e1ea988dd05")

	filter := bson.M{"_id": objID}
	// update := bson.M{"$set": bson.M{"isVerified": true}}
	database, err := db.NewMongoDB()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	res, err := database.UsersCollectionName.FindOne(context.Background(), filter).DecodeBytes()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Println("result: ", res)
}
