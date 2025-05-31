package usershandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/romanpitatelev/denet/internal/controller/rest/common"
	"github.com/romanpitatelev/denet/internal/entity"
)

type usersService interface {
	GetUser(ctx context.Context, userID entity.UserID) (entity.User, error)
	UpdateUser(ctx context.Context, userID entity.UserID, updatedUser entity.UserUpdate) (entity.User, error)
	DeleteUser(ctx context.Context, userID entity.UserID) error
	GetUsers(ctx context.Context, request entity.ListRequest) ([]entity.User, error)
}

type Handler struct {
	usersService usersService
}

func New(usersService usersService) *Handler {
	return &Handler{
		usersService: usersService,
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	userInfo, err := h.usersService.GetUser(ctx, entity.UserID(userID))
	if err != nil {
		common.ErrorResponse(w, "error getting user", err)

		return
	}

	common.OkResponse(w, http.StatusOK, userInfo)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	var user entity.UserUpdate

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)

		return
	}

	updatedUser, err := h.usersService.UpdateUser(ctx, entity.UserID(userID), user)
	if err != nil {
		common.ErrorResponse(w, "error updating user", err)

		return
	}

	common.OkResponse(w, http.StatusOK, updatedUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	if err = h.usersService.DeleteUser(ctx, entity.UserID(userID)); err != nil {
		common.ErrorResponse(w, "error deleting user", err)

		return
	}

	common.OkResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	request := common.GetListRequest(r)

	users, err := h.usersService.GetUsers(r.Context(), request)
	if err != nil {
		common.ErrorResponse(w, "error listing users", err)

		return
	}

	common.OkResponse(w, http.StatusOK, users)
}
