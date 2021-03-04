package model

import (
	"time"
)

// User type represents one application user.
type User struct {
	ID         int
	Name       string
	Email      string
	Password   string
	Superadmin bool
	Active     *time.Time
	Created    *time.Time
	Roles      *[]UserRole
	Links      *[]UserLink
}

// Role type represents one user's role.
type Role struct {
	ID      int
	Name    string
	Active  *time.Time
	Created *time.Time
	Users   *[]UserRole
	Links   *[]RoleLink
}

// UserRole type represents one connection between user and role types.
type UserRole struct {
	ID   int
	User *User
	Role *Role
}

// UserLink type represents one connection between user and link types.
type UserLink struct {
	ID    int
	User  *User
	Link  *Link
	Owner bool
}

// RoleLink type represents one connection between role and link types.
type RoleLink struct {
	ID   int
	Role *Role
	Link *Link
}
