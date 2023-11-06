package collection

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DatasetCollectionName CollectionName = "dataset"

type Dataset struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Descriptor [128]float32       `bson:"Descriptor"`
	UserID     string             `bson:"UserID"`
}
