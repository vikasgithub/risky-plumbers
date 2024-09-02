package risk

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vikasgithub/risky-plumbers/internal/entity"
	errorstype "github.com/vikasgithub/risky-plumbers/internal/errors"
	"github.com/vikasgithub/risky-plumbers/internal/log"
	"net/http"
	"strconv"
)

type resource struct {
	service Service
	logger  log.Logger
}

type RiskResponse struct {
	*entity.Risk
}

func (rr *RiskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewRiskResponse(risk *entity.Risk) *RiskResponse {
	return &RiskResponse{Risk: risk}
}

func NewRiskListResponse(articles []*entity.Risk) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, NewRiskResponse(article))
	}
	return list
}

func RegisterHandlers(r *chi.Mux, service Service) {
	res := resource{service, log.New()}

	r.Get("/risks/{id}", res.get)
	r.Get("/risks", res.getAll)
	r.Post("/risks", res.post)
}

func (res resource) get(w http.ResponseWriter, r *http.Request) {
	risk, err := res.service.Get(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, errorstype.ErrRecordNotFound) {
			render.Render(w, r, errorstype.ErrResponseNotFound)
		} else {
			render.Render(w, r, errorstype.ErrInvalidRequest(err))
		}
		return
	}

	render.Render(w, r, NewRiskResponse(risk))
}

// TODO need to implement paging and enhance the response with offset and limit
func (res resource) getAll(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 100
	if r.URL.Query().Get("offset") != "" {
		offsetParam, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			render.Render(w, r,
				errorstype.ErrInvalidRequest(
					fmt.Errorf("invalid offset: %s", r.URL.Query().Get("offset"))))
			return
		}
		offset = offsetParam
	}
	if r.URL.Query().Get("limit") != "" {
		limitParam, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			render.Render(w, r,
				errorstype.ErrInvalidRequest(
					fmt.Errorf("invalid limit: %s", r.URL.Query().Get("limit"))))
			return
		}
		limit = limitParam
	}

	res.logger.Infof("offset: %d, limit: %d", offset, limit)
	risks, err := res.service.GetAll(r.Context(), offset, limit)
	if err != nil {
		render.Render(w, r, errorstype.ErrInvalidRequest(err))
		return
	}
	render.RenderList(w, r, NewRiskListResponse(risks))
}

func (res resource) post(w http.ResponseWriter, r *http.Request) {
	createRequest := &CreateRiskRequest{}
	if err := render.Bind(r, createRequest); err != nil {
		render.Render(w, r, errorstype.ErrInvalidRequest(err))
		return
	}

	risk, err := res.service.Create(r.Context(), createRequest)
	if err != nil {
		render.Render(w, r, errorstype.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, NewRiskResponse(risk))
	if err != nil {
		render.Render(w, r, errorstype.ErrRender(err))
	}
}
