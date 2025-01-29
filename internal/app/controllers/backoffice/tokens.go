package backoffice

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type TokensController interface {
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type tokensController struct {
	tokens services.Tokens
}

func NewTokensController(tokens services.Tokens) TokensController {
	return &tokensController{
		tokens: tokens,
	}
}

//nolint:dupl
func (c *tokensController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paginator := services.NewPagination(r)
	rows, total, err := c.tokens.List(r.Context(), paginator)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	collection := make([]serializers.TokenSerializer, 0, len(rows))

	for _, row := range rows {
		collection = append(collection, serializers.TokenSerializer{
			ID:        row.ID,
			UserId:    row.UserId,
			Type:      row.Type,
			Value:     row.Value,
			ExpiresAt: row.ExpiresAt,
		})
	}

	response := serializers.PaginationResponse[serializers.TokenSerializer]{
		Data: collection,
		Meta: serializers.PaginationMeta{
			Page:  int(paginator.Page),
			Per:   int(paginator.PerPage),
			Total: total,
		},
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *tokensController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	_, err := c.tokens.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
