package models

import "time"

type Note struct {
	ID            int       `json:"id"`
	OwnerID       int       `json:"ownerId"`
	OwnerUsername string    `json:"ownerUsername,omitempty"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Shared        bool      `json:"shared"`
	Favorite      bool      `json:"favorite"`
	Updated       time.Time `json:"updatedAt"`
}
