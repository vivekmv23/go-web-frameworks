package database

import (
	"log"

	"github.com/google/uuid"
	"github.com/vivekmv23/go-web-frameworks/lib"
)

func TestDbClient() {

	id1 := uuid.New()

	itemToSave := lib.Item{
		Id:          id1,
		Name:        "some name",
		Value:       123,
		Description: "some item description",
		Active:      true,
	}

	if err := SaveItem(&itemToSave); err != nil {
		log.Printf("failed to save item: %s\n", err)
	} else {
		log.Println("saved item with id:", id1)
	}

	foundItem, err := GetItemById(id1)
	if err != nil {
		log.Println("failed to get:", err)
	} else {
		log.Println("found item with id:", foundItem.Id)
	}

	itemToSave.Value = 321

	updatedItem, err := UpdateItem(itemToSave, foundItem.UpdatedOn.String())

	if err != nil {
		log.Println("failed to update:", err)
	} else {
		log.Println("updated item with value:", updatedItem.Value)
	}

	items, err := GetAllItems()

	if err != nil {
		log.Println("failed to get all:", err)
	} else {
		log.Printf("found %d items", len(items))
		for idx, itm := range items {
			log.Printf("Items %d: Value: %d", idx, itm.Value)
		}
	}

	if err := DeleteItemById(id1); err != nil {
		log.Println("failed to delete:", err)
	} else {
		log.Println("deleted item with id:", id1)
	}

	// Not found tests
	id2 := uuid.New()

	if _, err := GetItemById(id2); err != nil {
		log.Println("failed to get:", err)
	} else {
		log.Println("found item with id:", id2)
	}

	if err := DeleteItemById(id2); err != nil {
		log.Println("failed to delete:", err)
	} else {
		log.Println("deleted item with id:", id2)
	}

	itemToSave.Id = id2
	updatedItem, err = UpdateItem(itemToSave, foundItem.UpdatedOn.String())
	if err != nil {
		log.Println("failed to update:", err)
	} else {
		log.Println("updated item with value:", updatedItem.Value)
	}
}
