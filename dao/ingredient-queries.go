package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetIngredientsNotInList(ingredients []primitive.ObjectID) interface{} {
	return bson.D{{"_id", bson.D{{"$nin", ingredients}}}}
}
