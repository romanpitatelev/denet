package taskservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/romanpitatelev/denet/internal/entity"
)

type tasksStore interface {
	Task(ctx context.Context, userID entity.UserID, task entity.Task) (entity.TaskResponse, error)
	ReferralTask()
}

type Service struct {
	tasksStore tasksStore
}

func New(tasksStore tasksStore) *Service {
	return &Service{
		tasksStore: tasksStore,
	}
}

func (s *Service) Task(ctx context.Context, userID entity.UserID, task entity.Task) (entity.TaskResponse, error) {
	taskValidated, err := task.Validate()
	if err != nil {
		return entity.TaskResponse{}, fmt.Errorf("failed to validate task: %w", err)
	}

	taskValidated.ID = entity.TaskID(uuid.New())

	taskResponse, err := s.tasksStore.Task(ctx, userID, taskValidated)
	if err != nil {
		return entity.TaskResponse{}, fmt.Errorf("failed to complete task: %w", err)
	}

	return taskResponse, nil
}
