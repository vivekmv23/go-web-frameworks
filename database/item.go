package database

import (
	"github.com/google/uuid"
	"time"
)

type Item struct {
	Id          uuid.UUID `bson:"id"`
	Name        string    `bson:"nam"`
	Value       int       `bson:"val"`
	Description string    `bson:"dsc"`
	Active      bool      `bson:"act"`
	CreatedOn   time.Time `bson:"con"`
	UpdatedOn   time.Time `bson:"uon"`
}
