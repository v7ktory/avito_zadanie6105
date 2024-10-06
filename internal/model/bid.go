package model

import "github.com/google/uuid"

type Bid struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TenderID    uuid.UUID `json:"tenderId"`
	AuthorType  string    `json:"authorType"`
	AuthorID    uuid.UUID `json:"authorId"`
}
