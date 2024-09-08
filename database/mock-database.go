package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/vivekmv23/go-web-frameworks/lib"
)

var (
	i1 lib.Item = lib.Item{
		Id:        uuid.New(),
		Name:      "name 1",
		Value:     1,
		Active:    true,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}
	i2 lib.Item = lib.Item{
		Id:        uuid.New(),
		Name:      "name 2",
		Value:     2,
		Active:    true,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}
	i3 lib.Item = lib.Item{
		Id:        uuid.New(),
		Name:      "name 3",
		Value:     3,
		Active:    true,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
	}
)

// Should satisfy ItemDatabase interface
type MockedDataBase struct {
	err error
}

func NewMockedDatabase(err error) (*MockedDataBase) {
	return &MockedDataBase{err: err}
}

func (m *MockedDataBase) SaveItem(i *lib.Item) error {
	return m.err
}

func (m *MockedDataBase) GetItemById(id uuid.UUID) (lib.Item, error) {
	return i1, m.err
}

func (m *MockedDataBase) GetAllItems() ([]lib.Item, error) {
	items := make([]lib.Item, 3)
	items = append(items, i1, i2, i3)
	return items, m.err
}

func (m *MockedDataBase) DeleteItemById(id uuid.UUID) error {
	return m.err
}

func (m *MockedDataBase) UpdateItem(i lib.Item, ifMatch string) (lib.Item, error) {
	return i1, m.err
}
