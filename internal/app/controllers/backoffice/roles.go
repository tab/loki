package backoffice

import (
	"encoding/json"
	"net/http"

	"loki/internal/app/models/dto"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"loki/internal/app/errors"
	"loki/internal/app/models"
	"loki/internal/app/serializers"
	"loki/internal/app/services"
)

type RolesController interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type rolesController struct {
	roles services.Roles
}

func NewRolesController(roles services.Roles) RolesController {
	return &rolesController{
		roles: roles,
	}
}

//nolint:dupl
func (c *rolesController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paginator := services.NewPagination(r)
	rows, total, err := c.roles.List(r.Context(), paginator)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	collection := make([]serializers.RoleSerializer, 0, len(rows))

	for _, row := range rows {
		collection = append(collection, serializers.RoleSerializer{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
		})
	}

	response := serializers.PaginationResponse[serializers.RoleSerializer]{
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
func (c *rolesController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	record, err := c.roles.FindRoleDetailsById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.ErrRoleNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.RoleSerializer{
		ID:            record.ID,
		Name:          record.Name,
		Description:   record.Description,
		PermissionIDs: record.PermissionIDs,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *rolesController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.RoleRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.roles.Create(r.Context(), &models.Role{
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.RoleSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *rolesController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	var params dto.RoleRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.roles.Update(r.Context(), &models.Role{
		ID:            id,
		Name:          params.Name,
		Description:   params.Description,
		PermissionIDs: params.PermissionIDs,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.RoleSerializer{
		ID:          record.ID,
		Name:        record.Name,
		Description: record.Description,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *rolesController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	_, err := c.roles.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
