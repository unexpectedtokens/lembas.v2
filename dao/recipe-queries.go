package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Create non-generic recipe dao that uses generic dao under the hood
func IngredientsInListQuery(ids []primitive.ObjectID) interface{} {
	return bson.D{{"_id", bson.D{{"$in", ids}}}}
}
