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

type EmployeeRepo struct {
	*postgres.Postgres
}

func NewEmployeeRepo(pg *postgres.Postgres) *EmployeeRepo {
	return &EmployeeRepo{pg}
}

func (er *EmployeeRepo) GetByUsername(ctx context.Context, username string) (*entity.Employee, error) {
	sql, args, _ := er.Builder.
		Select("id, username, first_name, last_name, created_at, updated_at").
		From("employee").
		Where(squirrel.Eq{"username": username}).
		ToSql()

	var employee entity.Employee
	err := er.Pool.QueryRow(ctx, sql, args...).Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetUserByUsername - er.Pool.QueryRow: %v", err)
	}

	return &employee, nil
}

func (er *EmployeeRepo) GetByID(ctx context.Context, employeeID uuid.UUID) (*entity.Employee, error) {
	sql, args, _ := er.Builder.
		Select("id, username, first_name, last_name, created_at, updated_at").
		From("employee").
		Where(squirrel.Eq{"id": employeeID}).
		ToSql()

	var employee entity.Employee
	err := er.Pool.QueryRow(ctx, sql, args...).Scan(
		&employee.ID,
		&employee.Username,
		&employee.FirstName,
		&employee.LastName,
		&employee.CreatedAt,
		&employee.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrs.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetUserByID - er.Pool.QueryRow: %v", err)
	}

	return &employee, nil
}

func (er *EmployeeRepo) IsResponsible(ctx context.Context, employeeID uuid.UUID) (bool, error) {
	sql, args, _ := er.Builder.
		Select("id").
		From("organization_responsible").
		Where(squirrel.Eq{"user_id": employeeID}).
		ToSql()

	var id uuid.UUID
	err := er.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("UserRepo.IsResponsible - er.Pool.QueryRow: %v", err)
	}
	return true, nil
}
