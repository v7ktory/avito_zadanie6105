package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	ID              uuid.UUID
	Name            string
	Description     string
	ServiceType     string
	Status          string
	Version         int
	OrganizationID  uuid.UUID
	CreatorUsername string
	CreatedAt       time.Time
}
