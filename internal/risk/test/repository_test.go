package risktest

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vikasgithub/risky-plumbers/internal/entity"
	errorstype "github.com/vikasgithub/risky-plumbers/internal/errors"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	risk2 "github.com/vikasgithub/risky-plumbers/internal/risk"
	"testing"
)

func TestGetRecordFound(t *testing.T) {
	repo := risk2.NewRepository(log.New())
	repo.Create(context.Background(), &entity.Risk{
		ID:          "1",
		State:       "open",
		Title:       "title",
		Description: "desc",
	})
	risk, _ := repo.Get(context.Background(), "1")
	assert.NotEmpty(t, risk)
	assert.Equal(t, "1", risk.ID)
}

func TestGetRecordNotFound(t *testing.T) {
	repo := risk2.NewRepository(log.New())
	_, err := repo.Get(context.Background(), "1")
	assert.NotEmpty(t, err)
	assert.IsType(t, errorstype.ErrRecordNotFound, err)
}

func TestQueryRecordsFound(t *testing.T) {
	repo := risk2.NewRepository(log.New())
	repo.Create(context.Background(), &entity.Risk{ID: "1", State: "open", Title: "title", Description: "desc"})
	repo.Create(context.Background(), &entity.Risk{ID: "2", State: "open", Title: "title", Description: "desc"})
	repo.Create(context.Background(), &entity.Risk{ID: "3", State: "open", Title: "title", Description: "desc"})
	risks, _ := repo.Query(context.Background(), 0, 5)
	assert.NotEmpty(t, risks)
	assert.Equal(t, 3, len(risks))
}

func TestQueryRecordsNotFound(t *testing.T) {
	repo := risk2.NewRepository(log.New())
	risks, _ := repo.Query(context.Background(), 0, 5)
	assert.Empty(t, risks)
	assert.Equal(t, 0, len(risks))
}
