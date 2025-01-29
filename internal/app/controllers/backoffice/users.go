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

type UsersController interface {
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type usersController struct {
	users services.Users
}

func NewUsersController(users services.Users) UsersController {
	return &usersController{
		users: users,
	}
}

//nolint:dupl
func (c *usersController) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paginator := services.NewPagination(r)
	rows, total, err := c.users.List(r.Context(), paginator)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	collection := make([]serializers.UserSerializer, 0, len(rows))

	for _, row := range rows {
		collection = append(collection, serializers.UserSerializer{
			ID:             row.ID,
			IdentityNumber: row.IdentityNumber,
			PersonalCode:   row.PersonalCode,
			FirstName:      row.FirstName,
			LastName:       row.LastName,
		})
	}

	response := serializers.PaginationResponse[serializers.UserSerializer]{
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
func (c *usersController) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	record, err := c.users.FindUserDetailsById(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.UserSerializer{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
		RoleIDs:        record.RoleIDs,
		ScopeIDs:       record.ScopeIDs,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *usersController) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params dto.UserRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.users.Create(r.Context(), &models.User{
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.UserSerializer{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *usersController) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	var params dto.UserRequest
	if err := params.Validate(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	record, err := c.users.Update(r.Context(), &models.User{
		ID:             id,
		IdentityNumber: params.IdentityNumber,
		PersonalCode:   params.PersonalCode,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		RoleIDs:        params.RoleIDs,
		ScopeIDs:       params.ScopeIDs,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	response := serializers.UserSerializer{
		ID:             record.ID,
		IdentityNumber: record.IdentityNumber,
		PersonalCode:   record.PersonalCode,
		FirstName:      record.FirstName,
		LastName:       record.LastName,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

//nolint:dupl
func (c *usersController) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := uuid.MustParse(chi.URLParam(r, "id"))

	_, err := c.users.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
