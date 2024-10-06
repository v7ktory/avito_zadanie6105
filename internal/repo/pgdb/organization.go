package pgdb

import (
	"context"
	"errors"
	"fmt"

	"git.codenrock.com/tender/internal/entity"
	"git.codenrock.com/tender/internal/repo/repoerrs"
	"git.codenrock.com/tender/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type OrganizationRepo struct {
	*postgres.Postgres
}

func NewOrganizationRepo(pg *postgres.Postgres) *OrganizationRepo {
	return &OrganizationRepo{pg}
}

func (or *OrganizationRepo) GetOrganizationResponsible(ctx context.Context, organizationID uuid.UUID, employeeID uuid.UUID) (*entity.OrganizationResponsible, error) {
	sql, args, err := or.Builder.
		Select("id", "organization_id", "user_id").
		From("organization_responsible").
		Where(squirrel.And{
			squirrel.Eq{"organization_id": organizationID},
			squirrel.Eq{"user_id": employeeID},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepo.GetOrganizationResponsible - failed to build SQL query: %w", err)
	}

	var organizationResponsible entity.OrganizationResponsible
	err = or.Pool.QueryRow(ctx, sql, args...).Scan(
		&organizationResponsible.ID,
		&organizationResponsible.OrganizationID,
		&organizationResponsible.UserID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("OrganizationRepo.GetOrganizationResponsible - failed to execute query: %w", err)
	}

	return &organizationResponsible, nil
}
func (or *OrganizationRepo) IsResponsibleForTender(ctx context.Context, userID uuid.UUID, tenderID uuid.UUID) (bool, error) {
	sql, args, err := or.Builder.
		Select("count(*)").
		From("organization_responsible").
		Join("tender ON tender.organization_id = organization_responsible.organization_id").
		Where(squirrel.Eq{
			"organization_responsible.user_id": userID,
			"tender.id":                        tenderID,
		}).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("OrganizationRepo.IsResponsibleForTender - failed to build SQL query: %w", err)
	}

	var count int
	err = or.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("OrganizationRepo.IsResponsibleForTender - query execution failed: %w", err)
	}

	return count > 0, nil
}
