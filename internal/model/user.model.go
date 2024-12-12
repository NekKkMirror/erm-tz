package model

import "time"

// User structure of the user table
type User struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"create_at"`
}
