package risktest

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vikasgithub/risky-plumbers/internal/entity"
	errorstype "github.com/vikasgithub/risky-plumbers/internal/errors"
	"github.com/vikasgithub/risky-plumbers/internal/risk"
	mocks "github.com/vikasgithub/risky-plumbers/internal/risk/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	router := chi.NewRouter()
	riskService := &mocks.Service{}
	risk.RegisterHandlers(router, riskService)

	t.Run("Risk Not Found", func(t *testing.T) {
		riskService.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, errorstype.ErrRecordNotFound).Once()
		rq, _ := http.NewRequest("GET", "/risks/1", nil)
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusNotFound, rs.Result().StatusCode)
		fmt.Println(rs.Body.String())
		assert.Equal(t, `{"status":"Resource not found."}`, strings.Trim(rs.Body.String(), "\n"))
	})

	t.Run("Risk Found", func(t *testing.T) {
		riskEntity := &entity.Risk{ID: "2", State: "o", Title: "t", Description: "d"}
		riskService.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(riskEntity, nil).Once()
		rq, _ := http.NewRequest("GET", "/risks/2", nil)
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusOK, rs.Result().StatusCode)
		assert.Equal(t, `{"id":"2","state":"o","title":"t","description":"d"}`, strings.Trim(rs.Body.String(), "\n"))

		fmt.Println(rs.Body.String())
	})
}

func TestGetAll(t *testing.T) {
	router := chi.NewRouter()
	riskService := &mocks.Service{}
	risk.RegisterHandlers(router, riskService)

	t.Run("Invalid Offset", func(t *testing.T) {
		rq, _ := http.NewRequest("GET", "/risks?offset=a", nil)
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusBadRequest, rs.Result().StatusCode)
		assert.Equal(t, `{"status":"Invalid request.","error":"invalid offset: a"}`, strings.Trim(rs.Body.String(), "\n"))
	})

	t.Run("Invalid Limit", func(t *testing.T) {
		rq, _ := http.NewRequest("GET", "/risks?limit=a", nil)
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusBadRequest, rs.Result().StatusCode)
		assert.Equal(t, `{"status":"Invalid request.","error":"invalid limit: a"}`, strings.Trim(rs.Body.String(), "\n"))
	})

	t.Run("Test Error", func(t *testing.T) {
		riskService.On("GetAll",
			mock.Anything,
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(nil, errors.New("get error")).Once()
		rq, _ := http.NewRequest("GET", "/risks", nil)
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusBadRequest, rs.Result().StatusCode)
		assert.Equal(t, `{"status":"Invalid request.","error":"get error"}`, strings.Trim(rs.Body.String(), "\n"))
	})

	t.Run("Test Success", func(t *testing.T) {
		risks := []*entity.Risk{
			{ID: "1", State: "o", Title: "t", Description: "d"},
			{ID: "2", State: "o", Title: "t", Description: "d"},
		}
		riskService.On("GetAll",
			mock.Anything,
			mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).
			Return(risks, nil).Once()
		rq, _ := http.NewRequest("GET", "/risks", nil)
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusOK, rs.Result().StatusCode)
		assert.Equal(t,
			`[{"id":"1","state":"o","title":"t","description":"d"},{"id":"2","state":"o","title":"t","description":"d"}]`,
			strings.Trim(rs.Body.String(),
				"\n"))
	})
}

func TestCreate(t *testing.T) {
	router := chi.NewRouter()
	riskService := &mocks.Service{}
	risk.RegisterHandlers(router, riskService)

	t.Run("Invalid Risk Request Parameters", func(t *testing.T) {
		riskService.On("Create", mock.Anything, mock.Anything).
			Return(nil, errors.New("invalid request")).Once()
		rq, _ := http.NewRequest("POST", "/risks",
			bytes.NewBufferString(`{"state":"o","title":"t","description":"d"}`))
		rq.Header.Set("Content-Type", "application/json")
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusBadRequest, rs.Result().StatusCode)
		assert.Equal(t, `{"status":"Invalid request.","error":"invalid request"}`, strings.Trim(rs.Body.String(), "\n"))
	})

	t.Run("Test Success", func(t *testing.T) {
		riskEntity := &entity.Risk{ID: "2", State: "o", Title: "t", Description: "d"}
		riskService.On("Create", mock.Anything, mock.Anything).
			Return(riskEntity, nil).Once()
		rq, _ := http.NewRequest("POST", "/risks",
			bytes.NewBufferString(`{"state":"o","title":"t","description":"d"}`))
		rq.Header.Set("Content-Type", "application/json")
		rs := httptest.NewRecorder()
		router.ServeHTTP(rs, rq)
		assert.Equal(t, http.StatusCreated, rs.Result().StatusCode)
		assert.Equal(t, `{"id":"2","state":"o","title":"t","description":"d"}`, strings.Trim(rs.Body.String(), "\n"))
	})
}
