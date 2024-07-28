package database

import (
	"log"

	"github.com/google/uuid"
)

func TestDbClient() {

	id1 := uuid.New()

	itemToSave := Item{
		Id:          id1,
		Name:        "some name",
		Value:       123,
		Description: "some item description",
		Active:      true,
	}

	if err := Save(&itemToSave); err != nil {
		log.Printf("failed to save item: %s\n", err)
	} else {
		log.Println("saved item with id:", id1)
	}

	if _, err := GetById(id1); err != nil {
		log.Println("failed to get:", err)
	} else {
		log.Println("found item with id:", id1)
	}

	items, err := GetAll()

	if err != nil {
		log.Println("failed to get all:", err)
	} else {
		log.Printf("found %d items", len(items))
		for idx, itm := range items {
			log.Printf("Items %d: %#v", idx, itm)
		}
	}

	if err := DeleteById(id1); err != nil {
		log.Println("failed to delete:", err)
	} else {
		log.Println("deleted item with id:", id1)
	}

	// Not found tests
	id2 := uuid.New()

	if _, err := GetById(id2); err != nil {
		log.Println("failed to get:", err)
	} else {
		log.Println("found item with id:", id2)
	}

	if err := DeleteById(id2); err != nil {
		log.Println("failed to delete:", err)
	} else {
		log.Println("deleted item with id:", id2)
	}
}
