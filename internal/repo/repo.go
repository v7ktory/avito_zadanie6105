package repo

import (
	"context"

	"git.codenrock.com/tender/internal/entity"
	"git.codenrock.com/tender/internal/repo/pgdb"
	"git.codenrock.com/tender/pkg/postgres"
	"github.com/google/uuid"
)

type Employee interface {
	GetByUsername(ctx context.Context, username string) (*entity.Employee, error)
	GetByID(ctx context.Context, employeeID uuid.UUID) (*entity.Employee, error)
	IsResponsible(ctx context.Context, employeeID uuid.UUID) (bool, error)
}
type Tender interface {
	CreateTender(ctx context.Context, t *entity.Tender) (*entity.Tender, error)
	GetTenders(ctx context.Context, limit, offset int, serviceTypes []string) ([]*entity.Tender, error)
	GetUserTenders(ctx context.Context, limit, offset int, username string) ([]*entity.Tender, error)
	GetTenderByID(ctx context.Context, tenderID uuid.UUID) (*entity.Tender, error)
	UpdateTenderStatus(ctx context.Context, tenderID uuid.UUID, status string) (*entity.Tender, error)
	UpdateTender(ctx context.Context, tenderID uuid.UUID, updates map[string]interface{}) (*entity.Tender, error)
}

type Bid interface {
	CreateBid(ctx context.Context, b *entity.Bid) (*entity.Bid, error)
	FindByTenderAndAuthor(ctx context.Context, tenderID, authorID uuid.UUID) (*entity.Bid, error)
	GetUserBids(ctx context.Context, limit, offset int, authorID uuid.UUID) ([]*entity.Bid, error)
	GetBidsByTender(ctx context.Context, limit int, offset int, tenderID uuid.UUID) ([]*entity.Bid, error)
	GetBidByID(ctx context.Context, bidID uuid.UUID) (*entity.Bid, error)
	UpdateBidStatus(ctx context.Context, bidID uuid.UUID, status string) (*entity.Bid, error)
	UpdateBid(ctx context.Context, bidID uuid.UUID, updates map[string]interface{}) (*entity.Bid, error)
	UpdateBidDecision(ctx context.Context, bidID uuid.UUID, decision string) (*entity.Bid, error)
	UpdateBidFeedback(ctx context.Context, bidID uuid.UUID, feedback string) (*entity.Bid, error)
}
type Organization interface {
	GetOrganizationResponsible(ctx context.Context, organizationID uuid.UUID, employeeID uuid.UUID) (*entity.OrganizationResponsible, error)
	IsResponsibleForTender(ctx context.Context, userID uuid.UUID, tenderID uuid.UUID) (bool, error)
}
type Repositories struct {
	Tender
	Employee
	Organization
	Bid
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Tender:       pgdb.NewTenderRepo(pg),
		Employee:     pgdb.NewEmployeeRepo(pg),
		Organization: pgdb.NewOrganizationRepo(pg),
		Bid:          pgdb.NewBidRepo(pg),
	}
}
