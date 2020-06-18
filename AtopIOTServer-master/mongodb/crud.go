package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) Count(c string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	return collection.CountDocuments(ctx, bson.M{})
}

// FindOne find first result in collection c with query
func (d *DB) FindOne(c string, mongoquery bson.M) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	var result bson.M
	err := collection.FindOne(ctx, mongoquery).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Read read data form collection c.
func (d *DB) Read(c string, query bson.M) ([]bson.M, error) {
	pagePara := GetPageParameter(query)
	mongoQuery := GetMongoQuery(query)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	opt := options.Find()
	opt.Limit = &pagePara.Limit
	opt.Skip = &pagePara.Offset
	cursor, err := collection.Find(ctx, mongoQuery, opt)
	if err != nil {
		log.Println("READ find error", err)
		return nil, err
	}
	var data []bson.M
	err = cursor.All(context.TODO(), &data)
	if err != nil {
		log.Println("READ decode error", err)
		return nil, err
	}
	return data, nil
}

//UpdateID update exist document
func (d *DB) UpdateID(c, idHex string, data bson.M) (bson.M, error) {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	res := collection.FindOneAndUpdate(ctx,
		bson.M{"_id": id},
		bson.M{"$set": data},
		options.FindOneAndUpdate().SetUpsert(true))
	fmt.Println("update data", res)
	var doc = bson.M{}
	decodeErr := res.Decode(doc)
	return doc, decodeErr
}
