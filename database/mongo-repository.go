package database

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ITEM_DB         = "itemDB"
	ITEM_COLLECTION = "items"
)

func getMongoCollection() *mongo.Collection {

	// Production ready application should ideally form the URI with credentials from ENV variables
	client, err := mongo.Connect(context.TODO(), options.Client(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Printf("ERROR: failed to get client: %s", err)
	}

	return client.Database(ITEM_DB).Collection(ITEM_COLLECTION)
}

func Save(i *Item) error {
	mc := getMongoCollection()
	updateTimes(i)
	_, err := mc.InsertOne(context.TODO(), i)

	return mapDbError(err)

}

func GetById(id uuid.UUID) (Item, error) {
	mc := getMongoCollection()
	var i Item
	err := mc.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&i)
	return i, mapDbError(err, id)

}

func DeleteById(id uuid.UUID) error {
	mc := getMongoCollection()

	res, err := mc.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if res.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
	}

	return mapDbError(err, id)
}

func mapDbError(err error, arg ...any) error {

	if err == nil {
		return nil
	}

	if err == mongo.ErrNoDocuments {
		return &NotFound{id: arg[0]}
	}

	if mongo.IsDuplicateKeyError(err) {
		return &Conflict{}
	}

	return &Unclassified{err: err}
}

func updateTimes(i *Item) {
	now := time.Now()
	if i.CreatedOn.IsZero() {
		i.CreatedOn = now
	}
	i.UpdatedOn = now
}
