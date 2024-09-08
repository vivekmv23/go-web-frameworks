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

type ItemDatabase interface {
	SaveItem(i *lib.Item) error
	GetItemById(id uuid.UUID) (lib.Item, error)
	GetAllItems() ([]lib.Item, error)
	DeleteItemById(id uuid.UUID) error
	UpdateItem(i lib.Item, ifMatch string) (lib.Item, error)
}

type Database struct {
	connection_url string
}

func NewDatabase() ItemDatabase {
	return &Database{
		connection_url: "mongodb://localhost:27017",
	}
}

func NewDatabaseWithUrl(connection_url string) ItemDatabase {
	return &Database{
		connection_url: connection_url,
	}
}

func (d *Database) getMongoCollection() *mongo.Collection {

	// Production ready application should ideally form the URI with credentials from ENV variables
	client, err := mongo.Connect(context.TODO(), options.Client(), options.Client().ApplyURI(d.connection_url))

	if err != nil {
		log.Printf("ERROR: failed to get mongo client: %s", err)
	}

	return client.Database(ITEM_DB).Collection(ITEM_COLLECTION)
}

func (d *Database) SaveItem(i *lib.Item) error {
	mc := d.getMongoCollection()
	determinations(i)
	_, err := mc.InsertOne(context.TODO(), i)

	return mapDbError(err)

}

func (d *Database) GetItemById(id uuid.UUID) (lib.Item, error) {
	mc := d.getMongoCollection()
	var i lib.Item
	err := mc.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&i)
	return i, mapDbError(err, id)

}

func (d *Database) GetAllItems() ([]lib.Item, error) {

	mc := d.getMongoCollection()

	var i []lib.Item

	cur, err := mc.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return i, mapDbError(err)
	}

	err = cur.All(context.TODO(), &i)

	if i == nil {
		i = make([]lib.Item, 0)
	}

	return i, mapDbError(err)

}

func (d *Database) DeleteItemById(id uuid.UUID) error {
	mc := d.getMongoCollection()

	res, err := mc.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if res.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
	}

	return mapDbError(err, id)
}

func (d *Database) UpdateItem(i lib.Item, ifMatch string) (lib.Item, error) {
	var existingItem lib.Item

	existingItem, err := d.GetItemById(i.Id)

	if err != nil {
		return i, err
	}

	if existingItem.UpdatedOn.String() != ifMatch {
		return i, &Outdated{}
	}

	i.DbId = existingItem.DbId
	i.CreatedOn = existingItem.CreatedOn
	determinations(&i)

	iDoc, err := toDoc(i)
	if err != nil {
		return i, mapDbError(err)
	}

	mc := d.getMongoCollection()

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
		return &NotFound{Id: arg[0]}
	}

	if mongo.IsDuplicateKeyError(err) {
		return &Conflict{}
	}

	return &Unclassified{Err: err}
}

func determinations(i *lib.Item) {

	if i.Id == uuid.Nil {
		i.Id = uuid.New()
	}

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
