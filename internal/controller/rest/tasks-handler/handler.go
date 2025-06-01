package taskshandler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/romanpitatelev/denet/internal/controller/rest/common"
	"github.com/romanpitatelev/denet/internal/entity"
)

type taskService interface {
	Task(ctx context.Context, userID entity.UserID, task entity.Task) (entity.TaskResponse, error)
	ReferralTask(ctx context.Context, reference entity.Reference) (entity.ReferenceResponse, error)
}

type Handler struct {
	taskService taskService
}

func New(taskService taskService) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

func (h *Handler) Task(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	var task entity.Task

	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)

		return
	}

	taskResponse, err := h.taskService.Task(ctx, entity.UserID(userID), task)
	if err != nil {
		common.ErrorResponse(w, "error updating user", err)

		return
	}

	common.OkResponse(w, http.StatusOK, taskResponse)
}

func (h *Handler) ReferralTask(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	var reference entity.Reference

	if err = json.NewDecoder(r.Body).Decode(&reference); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)

		return
	}

	reference.UserID = entity.UserID(userID)

	referenceResponse, err := h.taskService.ReferralTask(ctx, reference)
	if err != nil {
		common.ErrorResponse(w, "error updating user", err)

		return
	}

	common.OkResponse(w, http.StatusOK, referenceResponse)
}
