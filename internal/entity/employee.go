package entity

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
