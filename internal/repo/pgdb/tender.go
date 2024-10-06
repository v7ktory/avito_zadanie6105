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
	"github.com/jackc/pgx/v5/pgconn"
)

type TenderRepo struct {
	*postgres.Postgres
}

func NewTenderRepo(pg *postgres.Postgres) *TenderRepo {
	return &TenderRepo{pg}
}

func (tr *TenderRepo) CreateTender(ctx context.Context, t *entity.Tender) (*entity.Tender, error) {
	sql, args, _ := tr.Builder.
		Insert("tender").
		Columns("id", "name", "description", "service_type", "status", "organization_id", "creator_username").
		Values(uuid.New(), t.Name, t.Description, t.ServiceType, "CREATED", t.OrganizationID, t.CreatorUsername).
		Suffix("RETURNING id, name, description, service_type, status,version, organization_id, creator_username, created_at").
		ToSql()

	var createdTender entity.Tender

	err := tr.Pool.QueryRow(ctx, sql, args...).Scan(
		&createdTender.ID,
		&createdTender.Name,
		&createdTender.Description,
		&createdTender.ServiceType,
		&createdTender.Status,
		&createdTender.Version,
		&createdTender.OrganizationID,
		&createdTender.CreatorUsername,
		&createdTender.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return nil, repoerrs.ErrAlreadyExists
			}
		}
		return nil, fmt.Errorf("pgdb - TenderRepo - CreateTender: %w", err)
	}

	return &createdTender, nil
}

func (tr *TenderRepo) GetTenders(ctx context.Context, limit int, offset int, serviceTypes []string) ([]*entity.Tender, error) {
	queryBuilder := tr.Builder.
		Select("id", "name", "description", "service_type", "status", "version", "organization_id", "creator_username", "created_at").
		From("tender")

	if len(serviceTypes) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"service_type": serviceTypes})
	}
	queryBuilder = queryBuilder.OrderBy("name ASC")

	sql, args, _ := queryBuilder.
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := tr.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Tender{}, nil
		}
		return nil, fmt.Errorf("pgdb - TenderRepo - GetTenders: %w", err)
	}
	defer rows.Close()

	var tenders []*entity.Tender
	for rows.Next() {
		var tender entity.Tender
		if err := rows.Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.Version,
			&tender.OrganizationID,
			&tender.CreatorUsername,
			&tender.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tenders = append(tenders, &tender)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tenders, nil
}

func (tr *TenderRepo) GetUserTenders(ctx context.Context, limit, offset int, username string) ([]*entity.Tender, error) {
	sql, args, _ := tr.Builder.Select("id", "name", "description", "service_type", "status", "version", "organization_id", "creator_username", "created_at").
		From("tender").
		Where(squirrel.Eq{"creator_username": username}).
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := tr.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Tender{}, nil
		}
		return nil, fmt.Errorf("pgdb - TenderRepo - GetUserTenders: %w", err)
	}
	defer rows.Close()

	var tenders []*entity.Tender
	for rows.Next() {
		var tender entity.Tender
		if err := rows.Scan(
			&tender.ID,
			&tender.Name,
			&tender.Description,
			&tender.ServiceType,
			&tender.Status,
			&tender.Version,
			&tender.OrganizationID,
			&tender.CreatorUsername,
			&tender.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tenders = append(tenders, &tender)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tenders, nil
}

func (tr *TenderRepo) GetTenderByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error) {
	sql, args, _ := tr.Builder.
		Select("id", "name", "description", "service_type", "status", "version", "organization_id", "creator_username", "created_at").
		From("tender").
		Where(squirrel.Eq{"id": tenderID}).
		ToSql()

	var tender entity.Tender
	err := tr.Pool.QueryRow(ctx, sql, args...).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("pgdb - TenderRepo - GetTenderByID: %w", err)
	}

	return &tender, nil
}

func (tr *TenderRepo) UpdateTenderStatus(ctx context.Context, tenderID uuid.UUID, status string) (*entity.Tender, error) {
	sql, args, _ := tr.Builder.
		Update("tender").
		Set("status", status).
		Set("version", squirrel.Expr("version + 1")).
		Where(squirrel.Eq{"id": tenderID}).
		Suffix("RETURNING id, name, description, service_type, status, version, organization_id, creator_username, created_at").
		ToSql()

	var tender entity.Tender
	err := tr.Pool.QueryRow(ctx, sql, args...).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("pgdb - TenderRepo - UpdateTenderStatus: %w", err)
	}

	return &tender, nil
}

func (tr *TenderRepo) UpdateTender(ctx context.Context, tenderID uuid.UUID, updates map[string]interface{}) (*entity.Tender, error) {
	sql, args, _ := tr.Builder.
		Update("tender").
		SetMap(updates).
		Where(squirrel.Eq{"id": tenderID}).
		Suffix("RETURNING id, name, description, service_type, status, version, organization_id, creator_username, created_at").
		ToSql()

	var tender entity.Tender
	err := tr.Pool.QueryRow(ctx, sql, args...).Scan(
		&tender.ID,
		&tender.Name,
		&tender.Description,
		&tender.ServiceType,
		&tender.Status,
		&tender.Version,
		&tender.OrganizationID,
		&tender.CreatorUsername,
		&tender.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("pgdb - TenderRepo - UpdateTender: %w", err)
	}

	return &tender, nil
}
