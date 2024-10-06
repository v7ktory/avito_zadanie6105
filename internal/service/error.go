package service

import "fmt"

var (
	ErrTenderAlreadyExists      = fmt.Errorf("tender already exists")
	ErrCannotCreateTender       = fmt.Errorf("cannot create tender")
	ErrCannotGetTenders         = fmt.Errorf("cannot get tenders")
	ErrTenderNotFound           = fmt.Errorf("tender not found")
	ErrCannotGetTender          = fmt.Errorf("cannot get tender")
	ErrEmployeeNotFound         = fmt.Errorf("employee not found")
	ErrCannotUpdateTenderStatus = fmt.Errorf("cannot update tender status")
	ErrCannotUpdateTender       = fmt.Errorf("cannot update tender")
	ErrCannotCreateBid          = fmt.Errorf("cannot create bid")
	ErrBidAlreadyExists         = fmt.Errorf("bid already exists")
	ErrBidsNotFound             = fmt.Errorf("bids not found")
	ErrCannotGetBids            = fmt.Errorf("cannot get bids")
	ErrBidNotFound              = fmt.Errorf("bid not found")
	ErrCannotGetBid             = fmt.Errorf("cannot get bid")
	ErrCannotUpdateBid          = fmt.Errorf("cannot update bid")
)
