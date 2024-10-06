package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID          uuid.UUID
	Name        string
	Description string
	Status      string
	TenderID    uuid.UUID
	AuthorType  string
	AuthorID    uuid.UUID
	Version     int
	Decision    string
	Feedback    string
	CreatedAt   time.Time
}
