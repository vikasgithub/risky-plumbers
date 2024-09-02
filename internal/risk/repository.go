package risk

import (
	"context"
	"github.com/vikasgithub/risky-plumbers/internal/entity"
	errorstype "github.com/vikasgithub/risky-plumbers/internal/errors"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	"sync"
)

type Repository interface {
	Get(ctx context.Context, id string) (*entity.Risk, error)
	Query(ctx context.Context, offset, limit int) ([]*entity.Risk, error)
	Create(ctx context.Context, risk *entity.Risk) error
}

// when connecting with real db, the following struct will contain db context
type repository struct {
	cache  sync.Map
	logger log.Logger
}

func (r *repository) Get(ctx context.Context, id string) (*entity.Risk, error) {
	value, ok := r.cache.Load(id)
	if !ok {
		return nil, errorstype.ErrRecordNotFound
	}
	return value.(*entity.Risk), nil
}

func (r *repository) Query(ctx context.Context, offset, limit int) ([]*entity.Risk, error) {
	var entities []*entity.Risk
	//TODO offset and limit is not implemented due to time constraints
	r.cache.Range(func(key, value interface{}) bool {
		entities = append(entities, value.(*entity.Risk))
		return true
	})
	return entities, nil
}

func (r *repository) Create(ctx context.Context, risk *entity.Risk) error {
	r.cache.Store(risk.ID, risk)

	//Store does not throw any error, therefore returning nil here for the error
	return nil
}

func NewRepository(logger log.Logger) Repository {
	return &repository{logger: logger}
}
