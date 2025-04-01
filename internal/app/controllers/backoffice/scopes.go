package backoffice

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/models/dto"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type ScopesController interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type scopesController struct {
	scopes services.Scopes
}

func NewScopesController(scopes services.Scopes) ScopesController {
	return &scopesController{
		scopes: scopes,
	}
}

//nolint:dupl
func (c *scopesController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paginator := services.NewPagination(r)
	rows, total, err := c.scopes.List(r.Context(), paginator)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	collection := make([]serializers.ScopeSerializer, 0, len(rows))

	for _, row := range rows {
		collection = append(collection, serializers.ScopeSerializer{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
		})
	}

	response := serializers.PaginationResponse[serializers.ScopeSerializer]{
		Data: collection,
		Meta: serializers.PaginationMeta{
			Page:  paginator.Page,
			Per:   paginator.PerPage,
			Total: total,
		},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *scopesController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	record, err := c.scopes.FindById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.ErrScopeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.ScopeSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *scopesController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.ScopeRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.scopes.Create(r.Context(), &models.Scope{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.ScopeSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *scopesController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	var params dto.ScopeRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.scopes.Update(r.Context(), &models.Scope{
		ID:          id,
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.ScopeSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *scopesController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	_, err := c.scopes.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
