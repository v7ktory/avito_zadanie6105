package service

import (
	"context"
	"fmt"
	sl "log/slog"

	"git.codenrock.com/tender/internal/repo"
	"github.com/google/uuid"
)

type EmployeeService struct {
	employeeRepo repo.Employee
}

func NewEmployeeService(employeeRepo repo.Employee) *EmployeeService {
	return &EmployeeService{
		employeeRepo: employeeRepo,
	}
}

func (es *EmployeeService) GetByUsername(ctx context.Context, username string) (*EmployeeOutput, error) {
	const op = "service - EmployeeService - GetByUsername"

	employee, err := es.employeeRepo.GetByUsername(ctx, username)
	if err != nil {
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &EmployeeOutput{
		ID:        employee.ID,
		Username:  employee.Username,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}, nil
}

func (es *EmployeeService) GetByID(ctx context.Context, employeeID uuid.UUID) (*EmployeeOutput, error) {
	const op = "service - EmployeeService - GetByID"

	employee, err := es.employeeRepo.GetByID(ctx, employeeID)
	if err != nil {
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &EmployeeOutput{
		ID:        employee.ID,
		Username:  employee.Username,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}, nil
}

func (es *EmployeeService) IsResponsible(ctx context.Context, employeeID uuid.UUID) (bool, error) {
	const op = "service - EmployeeService - IsResponsible"

	isResponsible, err := es.employeeRepo.IsResponsible(ctx, employeeID)
	if err != nil {
		sl.Error(op, sl.Any("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isResponsible, nil
}
