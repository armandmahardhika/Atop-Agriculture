package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Count count documents in c (collection)
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

// Read read data form collection c. (c,q)->[item]
func (d *DB) Read(c string, query bson.M) ([]bson.M, error) {
	pagePara := GetPageParameter(query)
	mongoQuery := GetMongoQuery(query)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	opt := options.Find()
	opt.SetLimit(pagePara.Limit)
	opt.SetSkip(pagePara.Offset)
	// opt.Limit = &pagePara.Limit
	// opt.Skip = &pagePara.Offset

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

//UpdateID update exist document (c,q)->(pdateItem
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

	var doc = bson.M{}
	decodeErr := res.Decode(doc)
	return doc, decodeErr
}

// Insert (c,data)->({id:xx}) id is inserted item's id
func (d *DB) Insert(c string, data bson.M) (bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)

	res, err := collection.InsertOne(ctx, data)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return bson.M{"id": res.InsertedID}, err
}

//checkExist check if query exist?
func (d *DB) checkExist(c string, query bson.M) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	fmt.Println("checkExist ", query)
	collection := d.db.Collection(c)
	cur, err := collection.Find(ctx, query)

	if err != nil {
		log.Println(err)

		return false
	}
	defer cur.Close(ctx)
	return cur.TryNext(ctx)
}

//CheckUniqueValue check the data's properties satisfy unique policy return true if pass
// example
// if pass, msg := db.CheckUniqueValue(c, body); !pass {
//		res.SendUniqueValueError(errors.New(msg))
//		return
// }
func (d *DB) CheckUniqueValue(c string, data bson.M) (bool, string) {

	s, err := d.getTableSettings(c)
	if err != nil {
		// can not read table setting , the unique testing should pass
		fmt.Println(err)
		return true, ""
	}
	fmt.Println("setting.unique ", s.Unique)
	for _, key := range s.Unique {
		if val, ok := data[key]; ok {
			fmt.Println("checking  ", c, key, val)
			if exist := d.checkExist(c, bson.M{key: val}); exist {
				fmt.Println("exist")
				return false, fmt.Sprintf("%s : %s exist", key, val)
			}
		}
	}
	return true, ""

}

// DeleteOne delect first document (c,q)->({count:xx})
func (d *DB) DeleteOne(c string, idHex string) (bson.M, error) {
	//mongoQuery := GetMongoQuery(query)
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return bson.M{}, err
	}
	return bson.M{"count": result.DeletedCount}, nil
}

// DeleteMany delete by input query  (c,q)->({count:xx})
func (d *DB) DeleteMany(c string, query bson.M) (bson.M, error) {
	mongoQuery := GetMongoQuery(query)
	if len(mongoQuery) == 0 {
		return bson.M{}, errors.New("Parameter q can not be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)

	deleteResult, err := collection.DeleteMany(ctx, mongoQuery)
	if err != nil {
		log.Println("DeleteMany error", err)
		return bson.M{}, err
	}

	return bson.M{"count": deleteResult.DeletedCount}, nil
}

// DeleteIDs delete by ids array
func (d *DB) DeleteIDs(c string, query bson.M) (bson.M, error) {
	idsHEX, ok := query["ids[]"]

	if !ok {
		return bson.M{}, errors.New("Parameter ids not found")
	}

	if len(idsHEX.([]string)) <= 0 {
		return bson.M{}, errors.New("Parameter ids length <= 0")
	}

	var ids []primitive.ObjectID
	for _, v := range idsHEX.([]string) {
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection(c)

	deleteResult, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		log.Println("DeleteIDs error", err)
		return bson.M{}, err
	}

	return bson.M{"count": deleteResult.DeletedCount}, nil
}
