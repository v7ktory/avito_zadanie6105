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

type BidRepo struct {
	*postgres.Postgres
}

func NewBidRepo(pg *postgres.Postgres) *BidRepo {
	return &BidRepo{
		pg,
	}
}
func (br *BidRepo) CreateBid(ctx context.Context, bid *entity.Bid) (*entity.Bid, error) {
	sql, args, _ := br.Builder.
		Insert("bid").
		Columns("id", "name", "description", "status", "tender_id", "author_type", "author_id").
		Values(uuid.New(), bid.Name, bid.Description, "CREATED", bid.TenderID, bid.AuthorType, bid.AuthorID).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	var createdBid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&createdBid.ID,
		&createdBid.Name,
		&createdBid.Description,
		&createdBid.Status,
		&createdBid.TenderID,
		&createdBid.AuthorType,
		&createdBid.AuthorID,
		&createdBid.Version,
		&createdBid.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return nil, repoerrs.ErrAlreadyExists
			}
		}
		return nil, fmt.Errorf("pgdb - BidRepo - CreateBid: %w", err)
	}

	return &createdBid, nil
}
func (br *BidRepo) FindByTenderAndAuthor(ctx context.Context, tenderID, authorID uuid.UUID) (*entity.Bid, error) {
	sql, args, _ := br.Builder.
		Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bids").
		Where(squirrel.Eq{"tender_id": tenderID, "author_id": authorID}).
		ToSql()

	var bid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&bid.ID, &bid.Name, &bid.Description, &bid.TenderID,
		&bid.AuthorType, &bid.AuthorID, &bid.Version, &bid.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, err
	}

	return &bid, nil
}
func (br *BidRepo) GetUserBids(ctx context.Context, limit, offset int, authorID uuid.UUID) ([]*entity.Bid, error) {
	sql, args, _ := br.Builder.Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bid").
		Where(squirrel.Eq{"author_id": authorID}).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy("name ASC").
		ToSql()

	rows, err := br.Pool.Query(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var bids []*entity.Bid
	for rows.Next() {
		var bid entity.Bid
		if err := rows.Scan(
			&bid.ID,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderID,
			&bid.AuthorType,
			&bid.AuthorID,
			&bid.Version,
			&bid.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		bids = append(bids, &bid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return bids, nil
}

func (br *BidRepo) GetBidsByTender(ctx context.Context, limit int, offset int, tenderID uuid.UUID) ([]*entity.Bid, error) {
	sql, args, _ := br.Builder.Select(
		"id",
		"name",
		"description",
		"status",
		"tender_id",
		"author_type",
		"author_id",
		"version",
		"created_at",
	).
		From("bid").
		Where(squirrel.Eq{"tender_id": tenderID}).
		OrderBy("name ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	rows, err := br.Pool.Query(ctx, sql, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var bids []*entity.Bid
	for rows.Next() {
		var bid entity.Bid
		if err := rows.Scan(
			&bid.ID,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.TenderID,
			&bid.AuthorType,
			&bid.AuthorID,
			&bid.Version,
			&bid.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		bids = append(bids, &bid)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return bids, nil
}

func (br *BidRepo) GetBidByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error) {

	sql, args, _ := br.Builder.Select("id", "name", "description", "status", "tender_id", "author_type", "author_id", "version", "created_at").
		From("bid").
		Where(squirrel.Eq{"id": bidID}).
		ToSql()

	var bid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return &bid, nil
}

func (br *BidRepo) UpdateBidStatus(ctx context.Context, bidID uuid.UUID, status string) (*entity.Bid, error) {
	sql, args, _ := br.Builder.
		Update("bid").
		Set("status", status).
		Set("version", squirrel.Expr("version + 1")).
		Where(squirrel.Eq{"id": bidID}).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	var bid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &bid, nil
}
func (br *BidRepo) UpdateBid(ctx context.Context, bidID uuid.UUID, updates map[string]interface{}) (*entity.Bid, error) {
	sql, args, _ := br.Builder.
		Update("bid").
		SetMap(updates).
		Set("version", squirrel.Expr("version + 1")).
		Where(squirrel.Eq{"id": bidID}).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	var bid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &bid, nil
}
func (br *BidRepo) UpdateBidDecision(ctx context.Context, bidID uuid.UUID, decision string) (*entity.Bid, error) {
	sql, args, _ := br.Builder.
		Update("bid").
		Set("decision", decision).
		Set("version", squirrel.Expr("version + 1")).
		Where(squirrel.Eq{"id": bidID}).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	var bid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &bid, nil
}

func (br *BidRepo) UpdateBidFeedback(ctx context.Context, bidID uuid.UUID, feedback string) (*entity.Bid, error) {
	sql, args, _ := br.Builder.
		Update("bid").
		Set("feedback", feedback).
		Set("version", squirrel.Expr("version + 1")).
		Where(squirrel.Eq{"id": bidID}).
		Suffix("RETURNING id, name, description, status, tender_id, author_type, author_id, version, created_at").
		ToSql()

	var bid entity.Bid
	err := br.Pool.QueryRow(ctx, sql, args...).Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.TenderID,
		&bid.AuthorType,
		&bid.AuthorID,
		&bid.Version,
		&bid.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &bid, nil
}
