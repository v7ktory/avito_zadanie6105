package service

import (
	"context"
	"errors"
	sl "log/slog"

	"git.codenrock.com/tender/internal/entity"
	"git.codenrock.com/tender/internal/repo"
	"git.codenrock.com/tender/internal/repo/repoerrs"
	"github.com/google/uuid"
)

type TenderService struct {
	tenderRepo repo.Tender
}

func NewTenderService(tenderRepo repo.Tender) *TenderService {
	return &TenderService{
		tenderRepo: tenderRepo,
	}
}
func (ts *TenderService) CreateTender(ctx context.Context, input *CreateTenderInput) (*TenderOutput, error) {
	const op = "service - TenderService - CreateTender"

	tender := &entity.Tender{
		Name:            input.Name,
		Description:     input.Description,
		ServiceType:     input.ServiceType,
		OrganizationID:  input.OrganizationID,
		CreatorUsername: input.CreatorUsername,
	}

	output, err := ts.tenderRepo.CreateTender(ctx, tender)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return nil, ErrTenderAlreadyExists
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotCreateTender
	}

	return &TenderOutput{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		Status:      output.Status,
		ServiceType: output.ServiceType,
		Version:     output.Version,
		CreatedAt:   output.CreatedAt,
	}, nil
}

func (ts *TenderService) GetTenders(ctx context.Context, input *GetTendersInput) ([]*TenderOutput, error) {
	const op = "service - TenderService - GetTenders"

	tenders, err := ts.tenderRepo.GetTenders(ctx, input.Limit, input.Offset, input.ServiceTypes)
	if err != nil {
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotGetTenders
	}

	result := make([]*TenderOutput, 0, len(tenders))
	for _, t := range tenders {
		result = append(result, &TenderOutput{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			ServiceType: t.ServiceType,
			Status:      t.Status,
			Version:     t.Version,
			CreatedAt:   t.CreatedAt,
		})
	}

	return result, nil
}

func (ts *TenderService) GetUserTenders(ctx context.Context, input *GetUserTendersInput) ([]*TenderOutput, error) {
	const op = "service - TenderService - GetUserTenders"

	tenders, err := ts.tenderRepo.GetUserTenders(ctx, input.Limit, input.Offset, input.Username)
	if err != nil {
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotGetTenders
	}

	result := make([]*TenderOutput, 0, len(tenders))
	for _, t := range tenders {
		result = append(result, &TenderOutput{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			ServiceType: t.ServiceType,
			Status:      t.Status,
			Version:     t.Version,
			CreatedAt:   t.CreatedAt,
		})
	}

	return result, nil
}

func (ts *TenderService) GetTenderByID(ctx context.Context, tenderID uuid.UUID) (*TenderOutput, error) {
	const op = "service - TenderService - GetTenderByID"

	tender, err := ts.tenderRepo.GetTenderByID(ctx, tenderID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			sl.Error(op, sl.Any("error", "Tender not found"))
			return nil, ErrTenderNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotGetTender
	}

	return &TenderOutput{
		ID:              tender.ID,
		Name:            tender.Name,
		Description:     tender.Description,
		ServiceType:     tender.ServiceType,
		Status:          tender.Status,
		Version:         tender.Version,
		CreatorUsername: tender.CreatorUsername,
		OrganizationID:  tender.OrganizationID,
		CreatedAt:       tender.CreatedAt,
	}, nil
}

func (ts *TenderService) UpdateTenderStatus(ctx context.Context, input *UpdateTenderStatusInput) (*TenderOutput, error) {
	const op = "service - TenderService - UpdateTenderStatus"

	tender, err := ts.tenderRepo.UpdateTenderStatus(ctx, input.TenderID, input.Status)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			sl.Error(op, sl.Any("error", "Tender not found"))
			return nil, ErrTenderNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotUpdateTenderStatus
	}

	return &TenderOutput{
		ID:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      tender.Status,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}, nil
}

func (ts *TenderService) UpdateTender(ctx context.Context, input *UpdateTenderInput) (*TenderOutput, error) {
	const op = "service - TenderService - UpdateTender"

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.ServiceType != "" {
		updates["service_Type"] = input.ServiceType
	}

	tender, err := ts.tenderRepo.UpdateTender(ctx, input.TenderID, updates)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			sl.Error(op, sl.Any("error", "Tender not found"))
			return nil, ErrTenderNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotUpdateTender
	}

	return &TenderOutput{
		ID:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description,
		ServiceType: tender.ServiceType,
		Status:      tender.Status,
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt,
	}, nil
}
