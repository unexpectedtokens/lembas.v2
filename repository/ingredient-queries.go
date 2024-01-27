package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetIngredientsNotInList(ingredients []primitive.ObjectID) interface{} {
	return bson.D{{Key: "_id", Value: bson.D{{Key: "$nin", Value: ingredients}}}}
}
