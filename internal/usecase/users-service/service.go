package usersservice

import (
	"context"
	"fmt"

	"github.com/romanpitatelev/denet/internal/entity"
)

type usersStore interface {
	GetUser(ctx context.Context, userID entity.UserID) (entity.User, error)
	UpdateUser(ctx context.Context, userID entity.UserID, updatedUser entity.UserUpdate) (entity.User, error)
	DeleteUser(ctx context.Context, userID entity.UserID) error
	GetUsers(ctx context.Context, request entity.ListRequest) ([]entity.User, error)
}

type Service struct {
	usersStore usersStore
}

func New(usersStore usersStore) *Service {
	return &Service{
		usersStore: usersStore,
	}
}

func (s *Service) GetUser(ctx context.Context, userID entity.UserID) (entity.User, error) {
	user, err := s.usersStore.GetUser(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID entity.UserID, newInfoUser entity.UserUpdate) (entity.User, error) {
	_, err := s.usersStore.GetUser(ctx, userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("user not found: %w", err)
	}

	newInfoUserValidated, err := newInfoUser.Validate()
	if err != nil {
		return entity.User{}, fmt.Errorf("new info validation failed: %w", err)
	}

	updatedUser, err := s.usersStore.UpdateUser(ctx, userID, newInfoUserValidated)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to update user info: %w", err)
	}

	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID entity.UserID) error {
	if err := s.usersStore.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *Service) GetUsers(ctx context.Context, request entity.ListRequest) ([]entity.User, error) {
	users, err := s.usersStore.GetUsers(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("GetUsers %w", err)
	}

	return users, nil
}
