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

type BidService struct {
	bidRepo repo.Bid
}

func NewBidService(bidRepo repo.Bid) *BidService {
	return &BidService{
		bidRepo: bidRepo,
	}
}

func (bs *BidService) CreateBid(ctx context.Context, input *CreateBidInput) (*BidOutput, error) {
	const op = "service - BidService - CreateBid"

	bid := &entity.Bid{
		Name:        input.Name,
		Description: input.Description,
		TenderID:    input.TenderID,
		AuthorType:  input.AuthorType,
		AuthorID:    input.AuthorID,
	}

	createdBid, err := bs.bidRepo.CreateBid(ctx, bid)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return nil, ErrBidAlreadyExists
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotCreateBid
	}

	return &BidOutput{
		ID:         createdBid.ID,
		Name:       createdBid.Name,
		Status:     createdBid.Status,
		AuthorType: createdBid.AuthorType,
		AuthorID:   createdBid.AuthorID,
		Version:    createdBid.Version,
		CreatedAt:  createdBid.CreatedAt,
	}, nil
}
func (bs *BidService) GetBidByTenderAndAuthor(ctx context.Context, tenderID uuid.UUID, authorID uuid.UUID) (*BidOutput, error) {
	bid, err := bs.bidRepo.FindByTenderAndAuthor(ctx, tenderID, authorID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		return nil, err
	}
	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}

func (bs *BidService) GetUserBids(ctx context.Context, input *GetBidsByUsernameInput) ([]*BidOutput, error) {
	const op = "service - BidService - GetUserBids"

	bids, err := bs.bidRepo.GetUserBids(ctx, input.Limit, input.Offset, input.AuthorID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidsNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotGetBids
	}

	var output []*BidOutput
	for _, bid := range bids {
		output = append(output, &BidOutput{
			ID:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			Version:    bid.Version,
			CreatedAt:  bid.CreatedAt,
		})
	}

	return output, nil
}

func (bs *BidService) GetBidsByTender(ctx context.Context, input *GetBidsByTenderInput) ([]*BidOutput, error) {
	const op = "service - BidService - GetBidsByTender"

	bids, err := bs.bidRepo.GetBidsByTender(ctx, input.Limit, input.Offset, input.TenderID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidsNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotGetBids
	}

	var output []*BidOutput
	for _, bid := range bids {
		output = append(output, &BidOutput{
			ID:         bid.ID,
			Name:       bid.Name,
			Status:     bid.Status,
			AuthorType: bid.AuthorType,
			AuthorID:   bid.AuthorID,
			Version:    bid.Version,
			CreatedAt:  bid.CreatedAt,
		})
	}

	return output, nil
}

func (bs *BidService) GetBidByID(ctx context.Context, bidID uuid.UUID) (*BidOutput, error) {
	const op = "service - BidService - GetBidByID"

	bid, err := bs.bidRepo.GetBidByID(ctx, bidID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotGetBid
	}

	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		TenderID:   bid.TenderID,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}

func (bs *BidService) UpdateBidStatus(ctx context.Context, input *UpdateBidStatusInput) (*BidOutput, error) {
	const op = "service - BidService - UpdateBidStatus"

	bid, err := bs.bidRepo.UpdateBidStatus(ctx, input.BidID, input.Status)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotUpdateBid
	}

	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}
func (bs *BidService) UpdateBid(ctx context.Context, input *UpdateBidInput) (*BidOutput, error) {
	const op = "service - BidService - UpdateBid"

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["Description"] = input.Description
	}

	bid, err := bs.bidRepo.UpdateBid(ctx, input.BidID, updates)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotUpdateBid
	}

	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}
func (bs *BidService) UpdateBidDecision(ctx context.Context, input *UpdateBidDecisionInput) (*BidOutput, error) {
	const op = "service - BidService - UpdateBidDecision"

	bid, err := bs.bidRepo.UpdateBidDecision(ctx, input.BidID, input.Decision)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotUpdateBid
	}

	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}

func (bs *BidService) UpdateBidFeedback(ctx context.Context, input *UpdateBidFeedbackInput) (*BidOutput, error) {
	const op = "service - BidService - UpdateBidFeedback"

	bid, err := bs.bidRepo.UpdateBidFeedback(ctx, input.BidID, input.Feedback)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return nil, ErrBidNotFound
		}
		sl.Error(op, sl.Any("error", err.Error()))
		return nil, ErrCannotUpdateBid
	}

	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.AuthorID,
		Version:    bid.Version,
		CreatedAt:  bid.CreatedAt,
	}, nil
}
