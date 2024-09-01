package lib

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	DbId        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Id          uuid.UUID          `bson:"id,omitempty" json:"id"`
	Name        string             `bson:"nam,omitempty" json:"name"`
	Value       int                `bson:"val,omitempty" json:"value"`
	Description string             `bson:"dsc,omitempty" json:"description"`
	Active      bool               `bson:"act,omitempty" json:"isActive"`
	CreatedOn   time.Time          `bson:"con,omitempty" json:"createdOn"`
	UpdatedOn   time.Time          `bson:"uon,omitempty" json:"updatedOn"`
}

type Error struct {
	Error string `json:"error"`
	Path  string `json:"path"`
}
