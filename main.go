package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/vivekmv23/go-web-frameworks/database"
)

func main() {
	testDbClient()
}

func testDbClient() {

	id1 := uuid.New()

	itemToSave := database.Item{
		Id:          id1,
		Name:        "some name",
		Value:       123,
		Description: "some item description",
		Active:      true,
	}

	if err := database.Save(&itemToSave); err != nil {
		log.Printf("failed to save item: %s\n", err)
	} else {
		log.Println("saved item with id:", id1)
	}

	if _, err := database.GetById(id1); err != nil {
		log.Println("failed to get:", err)
	} else {
		log.Println("found item with id:", id1)
	}

	if err := database.DeleteById(id1); err != nil {
		log.Println("failed to delete:", err)
	} else {
		log.Println("deleted item with id:", id1)
	}

	// Not found tests
	id2 := uuid.New()

	if _, err := database.GetById(id2); err != nil {
		log.Println("failed to get:", err)
	} else {
		log.Println("found item with id:", id2)
	}

	if err := database.DeleteById(id2); err != nil {
		log.Println("failed to delete:", err)
	} else {
		log.Println("deleted item with id:", id2)
	}
}
