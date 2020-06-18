package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserCollection collection name of users
var UserCollection = "users"

// type UserRequest struct {
// 	Name     string   `bson:"name"`
// 	Password string   `bson:"password"`
// 	Tags     []string `bson:"tags"`
// }

//ValidUser check name and password
func (d *DB) ValidUser(body bson.M) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tracecode := body["tracecode"]
	colleciont := d.db.Collection(UserCollection)
	var item bson.M
	if err := colleciont.FindOne(ctx, bson.M{"tracecode": tracecode}).Decode(&item); err != nil {
		//id := item["_id"].(primitive.ObjectID)
		//	return id.Hex(), err
		return "", err
	}
	if item["tracecode"] == body["tracecode"] {
		id := item["_id"].(primitive.ObjectID)
		return id.Hex(), nil
	}

	return "", fmt.Errorf("No Data Found")
}
