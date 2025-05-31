package taskhandler

import (
	"context"

	"github.com/romanpitatelev/denet/internal/entity"
)

type taskService interface {
	GetTopUsers(ctx context.Context, request entity.ListRequest) ([]entity.User, error)
	Task()
	ReferralTask()
}

type Handler struct {
	taskService taskService
}

func New(taskService taskService) *Handler {
	return &Handler{
		taskService: taskService,
	}
}
