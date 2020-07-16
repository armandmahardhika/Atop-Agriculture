package mongodb

import "go.mongodb.org/mongo-driver/bson"

func (d *DB) WriteMqttClientStatus(status bson.M) error
