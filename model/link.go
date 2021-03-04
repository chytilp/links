package model

import (
	"time"
)

// Category type represents one category saved in db.
type Category struct {
	ID       int
	Name     string
	ParentID int
	Active   *time.Time
	Created  *time.Time
}

// Link type represents one link object saved in db.
type Link struct {
	ID       int
	Link     string
	Name     string
	Category *Category
	Active   *time.Time
	Created  *time.Time
}
