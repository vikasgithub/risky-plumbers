package risk

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/vikasgithub/risky-plumbers/internal/entity"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	"net/http"
)

type Service interface {
	Get(ctx context.Context, id string) (*entity.Risk, error)
	GetAll(ctx context.Context, offset, limit int) ([]*entity.Risk, error)
	Create(ctx context.Context, input *CreateRiskRequest) (*entity.Risk, error)
}

type CreateRiskRequest struct {
	State       string `json:"state"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (cr *CreateRiskRequest) Bind(r *http.Request) error {
	return nil
}

func (cr *CreateRiskRequest) Validate() error {
	return validation.ValidateStruct(cr,
		validation.Field(&cr.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&cr.Description, validation.Required, validation.Length(0, 1024)),
		validation.Field(&cr.State, validation.Required,
			validation.In("open", "closed", "accepted", "investigating")),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func (s service) Get(ctx context.Context, id string) (*entity.Risk, error) {
	return s.repo.Get(ctx, id)
}

func (s service) GetAll(ctx context.Context, offset, limit int) ([]*entity.Risk, error) {
	return s.repo.Query(ctx, offset, limit)
}

func (s service) Create(ctx context.Context, input *CreateRiskRequest) (*entity.Risk, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	id := entity.GenerateID()
	err := s.repo.Create(ctx, &entity.Risk{
		ID:          id,
		State:       input.State,
		Title:       input.Title,
		Description: input.Description,
	})
	if err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, id)
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}
