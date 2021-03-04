package model

import (
	"time"
)

// Star type represents one link score (in stars 1-5) from particular user.
type Star struct {
	ID      int
	User    *User
	Link    *Link
	Stars   int
	Created *time.Time
}

// Note type represents one link note from particular user.
type Note struct {
	ID      int
	User    *User
	Link    *Link
	Note    string
	Private bool
	Created *time.Time
}
