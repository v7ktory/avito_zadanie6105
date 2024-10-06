package service

import (
	"context"
	"fmt"
	sl "log/slog"

	"git.codenrock.com/tender/internal/repo"
	"github.com/google/uuid"
)

type OrganizationService struct {
	organizationRepo repo.Organization
}

func NewOrganizationService(organizationRepo repo.Organization) *OrganizationService {
	return &OrganizationService{
		organizationRepo: organizationRepo,
	}
}

func (os *OrganizationService) GetOrganizationResponsible(ctx context.Context, input *OrganizationResponsibleInput) (*OrganizationResponsibleOutput, error) {
	const op = "service - OrganizationService - GetOrganizationResponsible"

	organizationResponsible, err := os.organizationRepo.GetOrganizationResponsible(ctx, input.OrganizationID, input.EmployeeID)
	if err != nil {
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &OrganizationResponsibleOutput{
		ID:             organizationResponsible.ID,
		OrganizationID: organizationResponsible.OrganizationID,
		EmployeeID:     organizationResponsible.UserID,
	}, nil
}

func (os *OrganizationService) IsResponsibleForTender(ctx context.Context, userID uuid.UUID, tenderID uuid.UUID) (bool, error) {
	return os.organizationRepo.IsResponsibleForTender(ctx, userID, tenderID)
}
