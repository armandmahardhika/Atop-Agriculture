package mongodb

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB struct of monge database
type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

var instance *DB
var once sync.Once

const defaultURL string = "mongodb://localhost:27017"

// GetDB get unique database pointer
func GetDB() *DB {
	once.Do(func() {
		instance = initDB(defaultURL)
	})
	return instance
}

// Initial database
func initDB(url string) *DB {
	const defaultDatabase = "atop"
	db := new(DB)
	var err error
	db.client, err = mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = db.client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = db.client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	db.db = db.client.Database(defaultDatabase)
	log.Printf("Database %s connected... \n", url)
	return db
}

func (d *DB) GetCollections() []string {
	result, _ := d.db.ListCollectionNames(context.TODO(), bson.M{})

	return result
}

func (d *DB) Test() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	collection := d.db.Collection("test")

	if _, err := collection.InsertOne(ctx, bson.M{"name": "hello"}); err != nil {
		log.Println("test insert fail")
	}

}
