package tasksservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/romanpitatelev/denet/internal/entity"
)

const referencePoint = 1

type tasksStore interface {
	Task(ctx context.Context, userID entity.UserID, task entity.Task) (entity.TaskResponse, error)
	ReferralTask(ctx context.Context, reference entity.Reference) (entity.ReferenceResponse, error)
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

func (s *Service) ReferralTask(ctx context.Context, reference entity.Reference) (entity.ReferenceResponse, error) {
	if err := reference.Validate(); err != nil {
		return entity.ReferenceResponse{}, fmt.Errorf("failed to validate reference: %w", err)
	}

	reference.ID = entity.ReferenceID(uuid.New())
	reference.Points = referencePoint

	referenceResponse, err := s.tasksStore.ReferralTask(ctx, reference)
	if err != nil {
		return entity.ReferenceResponse{}, fmt.Errorf("failed to complete reference task: %w", err)
	}

	return referenceResponse, nil
}
