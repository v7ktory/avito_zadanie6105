package model

import (
	"context"

	"github.com/google/uuid"
)

type Validator interface {
	Valid(ctx context.Context) map[string]string
}

func (t Tender) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if t.Name == "" {
		problems["name"] = "Name is required"
	} else if len([]rune(t.Name)) > 50 {
		problems["name"] = "Name cannot be longer than 50 characters"
	}

	if t.Description == "" {
		problems["description"] = "Description is required"
	}

	if t.ServiceType == "" {
		problems["serviceType"] = "Service type is required"
	} else if _, valid := validateServiceType(t.ServiceType); !valid {
		problems["serviceType"] = "Invalid service type"
	}

	return problems
}
func validateServiceType(serviceType string) (string, bool) {
	validServiceTypes := map[string]bool{
		"Construction": true,
		"Delivery":     true,
		"Manufacture":  true,
	}

	if valid, ok := validServiceTypes[serviceType]; ok && valid {
		return serviceType, true
	}
	return "", false
}

func (b *Bid) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if b.Name == "" {
		problems["name"] = "Name is required"
	} else if len([]rune(b.Name)) > 50 {
		problems["name"] = "Name cannot be longer than 50 characters"
	}

	if b.Description == "" {
		problems["description"] = "Description is required"
	}

	if b.TenderID == uuid.Nil {
		problems["tenderId"] = "Tender ID is required"
	}

	if b.AuthorType == "" {
		problems["authorType"] = "Author type is required"
	} else if !validateAuthorType(b.AuthorType) {
		problems["authorType"] = "Invalid author type"
	}

	if b.AuthorID == uuid.Nil {
		problems["authorId"] = "Author ID is required"
	}

	return problems
}

func validateAuthorType(authorType string) bool {
	validAuthorTypes := map[string]bool{
		"Organization": true,
		"User":         true,
	}

	return validAuthorTypes[authorType]
}
