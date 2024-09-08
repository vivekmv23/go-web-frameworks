package database

import (
	"fmt"
)

type NotFound struct {
	Id interface{}
}

func (n *NotFound) Error() string {
	return fmt.Sprintf("item with id %s not found", n.Id)
}

type Outdated struct {
}

func (o *Outdated) Error() string {
	return "update request is outdated"
}

type Conflict struct {
}

func (c *Conflict) Error() string {
	return "attempted to save item with same id"
}

type Unclassified struct {
	Err error
}

func (u *Unclassified) Error() string {
	return fmt.Sprintf("unclassified error: %s", u.Err)
}
