package database

import (
	"fmt"
)

type NotFound struct {
	id interface{}
}

func (n *NotFound) Error() string {
	return fmt.Sprintf("item with id %s not found", n.id)
}

type Conflict struct {
}

func (c *Conflict) Error() string {
	return "attempted to save item with same id"
}

type Unclassified struct {
	err error
}

func (u *Unclassified) Error() string {
	return fmt.Sprintf("unclassified error: %s", u.err)
}
