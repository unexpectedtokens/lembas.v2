package repository

import (
	"go.mongodb.org/mongo-driver/bson"
)

// TODO: Create non-generic recipe dao that uses generic dao under the hood
func IngredientsInListQuery(ids []string) interface{} {
	return bson.D{{"_id", bson.D{{"$in", ids}}}}
}
