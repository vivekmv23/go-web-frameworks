package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vivekmv23/go-web-frameworks/lib"
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
		log.Printf("ERROR: failed to get mongo client: %s", err)
	}

	return client.Database(ITEM_DB).Collection(ITEM_COLLECTION)
}

func SaveItem(i *lib.Item) error {
	mc := getMongoCollection()
	updateTimes(i)
	_, err := mc.InsertOne(context.TODO(), i)

	return mapDbError(err)

}

func GetItemById(id uuid.UUID) (lib.Item, error) {
	mc := getMongoCollection()
	var i lib.Item
	err := mc.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&i)
	return i, mapDbError(err, id)

}

func GetAllItems() ([]lib.Item, error) {

	mc := getMongoCollection()

	var i []lib.Item

	cur, err := mc.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return i, mapDbError(err)
	}

	err = cur.All(context.TODO(), &i)

	return i, mapDbError(err)

}

func DeleteItemById(id uuid.UUID) error {
	mc := getMongoCollection()

	res, err := mc.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if res.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
	}

	return mapDbError(err, id)
}

func UpdateItem(i lib.Item, ifMatch string) (lib.Item, error) {
	var existingItem lib.Item

	existingItem, err := GetItemById(i.Id)

	if err != nil {
		return i, err
	}

	if existingItem.UpdatedOn.String() != ifMatch {
		return i, fmt.Errorf("Update request is outdated: if-match(%s) != %s", ifMatch, existingItem.UpdatedOn.String())
	}

	i.DbId = existingItem.DbId
	i.CreatedOn = existingItem.CreatedOn
	updateTimes(&i)

	iDoc, err := toDoc(i)
	if err != nil {
		return i, mapDbError(err)
	}

	mc := getMongoCollection()

	filter := bson.D{{Key: "_id", Value: i.DbId}}
	update := bson.D{{Key: "$set", Value: iDoc}}

	res, err := mc.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return i, mapDbError(err)
	}

	if res.ModifiedCount == 1 {
		return i, nil
	}

	return i, fmt.Errorf("failed to update, updated count != 1")
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

func updateTimes(i *lib.Item) {
	now := time.Now()
	if i.CreatedOn.IsZero() {
		i.CreatedOn = now
	}
	i.UpdatedOn = now
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}
