package service

import (
	"context"
	"time"

	"git.codenrock.com/tender/internal/repo"
	"github.com/google/uuid"
)

type TenderOutput struct {
	ID              uuid.UUID
	Name            string
	Description     string
	ServiceType     string
	Status          string
	Version         int
	OrganizationID  uuid.UUID
	CreatorUsername string
	CreatedAt       time.Time
}
type CreateTenderInput struct {
	Name            string
	Description     string
	ServiceType     string
	OrganizationID  uuid.UUID
	CreatorUsername string
}

type GetTendersInput struct {
	Limit        int
	Offset       int
	ServiceTypes []string
}

type GetUserTendersInput struct {
	Limit    int
	Offset   int
	Username string
}

type GetTenderByIDOutput struct {
	ID              uuid.UUID
	Name            string
	Description     string
	ServiceType     string
	Status          string
	Version         int
	CreatorUsername string
	CreatedAt       time.Time
}
type UpdateTenderStatusInput struct {
	TenderID uuid.UUID
	Status   string
}

type UpdateTenderInput struct {
	TenderID    uuid.UUID
	Name        string
	Description string
	ServiceType string
}

type RollbackTenderInput struct {
	TenderID uuid.UUID
	Version  int
}
type Tender interface {
	CreateTender(ctx context.Context, input *CreateTenderInput) (*TenderOutput, error)
	GetTenders(ctx context.Context, input *GetTendersInput) ([]*TenderOutput, error)
	GetUserTenders(ctx context.Context, input *GetUserTendersInput) ([]*TenderOutput, error)
	GetTenderByID(ctx context.Context, tenderID uuid.UUID) (*TenderOutput, error)
	UpdateTenderStatus(ctx context.Context, input *UpdateTenderStatusInput) (*TenderOutput, error)
	UpdateTender(ctx context.Context, input *UpdateTenderInput) (*TenderOutput, error)
}

type EmployeeOutput struct {
	ID        uuid.UUID
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Employee interface {
	GetByUsername(ctx context.Context, username string) (*EmployeeOutput, error)
	GetByID(ctx context.Context, employeeID uuid.UUID) (*EmployeeOutput, error)
	IsResponsible(ctx context.Context, employeeID uuid.UUID) (bool, error)
}

type OrganizationResponsibleInput struct {
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

type OrganizationResponsibleOutput struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	EmployeeID     uuid.UUID
}

type Organization interface {
	GetOrganizationResponsible(ctx context.Context, input *OrganizationResponsibleInput) (*OrganizationResponsibleOutput, error)
	IsResponsibleForTender(ctx context.Context, userID uuid.UUID, tenderID uuid.UUID) (bool, error)
}
type BidOutput struct {
	ID         uuid.UUID
	Name       string
	Status     string
	TenderID   uuid.UUID
	AuthorType string
	AuthorID   uuid.UUID
	Version    int
	CreatedAt  time.Time
}

type CreateBidInput struct {
	Name        string
	Description string
	TenderID    uuid.UUID
	AuthorType  string
	AuthorID    uuid.UUID
}

type GetBidsByUsernameInput struct {
	Limit    int
	Offset   int
	AuthorID uuid.UUID
}

type GetBidsByTenderInput struct {
	Limit    int
	Offset   int
	TenderID uuid.UUID
}

type UpdateBidStatusInput struct {
	BidID  uuid.UUID
	Status string
}

type UpdateBidInput struct {
	BidID       uuid.UUID
	Name        string
	Description string
}
type UpdateBidDecisionInput struct {
	BidID    uuid.UUID
	Decision string
}
type UpdateBidFeedbackInput struct {
	BidID    uuid.UUID
	Feedback string
}
type Bid interface {
	CreateBid(ctx context.Context, input *CreateBidInput) (*BidOutput, error)
	GetBidByTenderAndAuthor(ctx context.Context, tenderID uuid.UUID, authorID uuid.UUID) (*BidOutput, error)
	GetUserBids(ctx context.Context, input *GetBidsByUsernameInput) ([]*BidOutput, error)
	GetBidsByTender(ctx context.Context, input *GetBidsByTenderInput) ([]*BidOutput, error)
	GetBidByID(ctx context.Context, bidID uuid.UUID) (*BidOutput, error)
	UpdateBidStatus(ctx context.Context, input *UpdateBidStatusInput) (*BidOutput, error)
	UpdateBid(ctx context.Context, input *UpdateBidInput) (*BidOutput, error)
	UpdateBidDecision(ctx context.Context, input *UpdateBidDecisionInput) (*BidOutput, error)
	UpdateBidFeedback(ctx context.Context, input *UpdateBidFeedbackInput) (*BidOutput, error)
}
type Services struct {
	Tender
	Employee
	Organization
	Bid
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Tender:       NewTenderService(deps.Repos.Tender),
		Employee:     NewEmployeeService(deps.Repos.Employee),
		Organization: NewOrganizationService(deps.Repos.Organization),
		Bid:          NewBidService(deps.Repos.Bid),
	}
}
