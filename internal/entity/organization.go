package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationType int

const (
	IE  OrganizationType = iota 
	LLC                        
	JSC                       
)

type Organization struct {
	ID          uuid.UUID
	Name        string
	Description string
	Type        OrganizationType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrganizationResponsible struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
}
