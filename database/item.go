package database

import (
	"github.com/google/uuid"
	"time"
)

type Item struct {
	Id          uuid.UUID `bson:"id" json:"id"`
	Name        string    `bson:"nam" json:"name"`
	Value       int       `bson:"val" json:"value"`
	Description string    `bson:"dsc" json:"description"`
	Active      bool      `bson:"act" json:"isActive"`
	CreatedOn   time.Time `bson:"con" json:"createdOn"`
	UpdatedOn   time.Time `bson:"uon" json:"updatedOn"`
}
