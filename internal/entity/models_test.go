package entity_test

import (
	"testing"

	"github.com/romanpitatelev/denet/internal/entity"
	"github.com/romanpitatelev/denet/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUserUpdateValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     entity.UserUpdate
		expect    entity.UserUpdate
		expectErr bool
	}{
		{
			name:      "valid email",
			input:     entity.UserUpdate{Email: utils.Pointer("some@email.ru")},
			expect:    entity.UserUpdate{Email: utils.Pointer("some@email.ru")},
			expectErr: false,
		},
		{
			name:      "invalid email",
			input:     entity.UserUpdate{Email: utils.Pointer("some")},
			expect:    entity.UserUpdate{},
			expectErr: true,
		},
		{
			name:      "no email - valid",
			input:     entity.UserUpdate{Name: utils.Pointer("name"), Role: utils.Pointer("role")},
			expect:    entity.UserUpdate{Name: utils.Pointer("name"), Role: utils.Pointer("role")},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answer, err := tt.input.Validate()
			if tt.expectErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expect, answer)
		})
	}
}

func TestTaskValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     entity.Task
		expect    entity.Task
		expectErr error
	}{
		{
			name:      "valid telegram task",
			input:     entity.Task{Type: "telegram"},
			expect:    entity.Task{Type: "telegram", Points: 3},
			expectErr: nil,
		},
		{
			name:      "valid telegram task",
			input:     entity.Task{Type: "telegram"},
			expect:    entity.Task{Type: "telegram", Points: 3},
			expectErr: nil,
		},
		{
			name:      "valid twitter task",
			input:     entity.Task{Type: "twitter"},
			expect:    entity.Task{Type: "twitter", Points: 2},
			expectErr: nil,
		},
		{
			name:      "invalid task type",
			input:     entity.Task{Type: "some"},
			expect:    entity.Task{},
			expectErr: entity.ErrInvalidTaskType,
		},
		{
			name:      "invalid points - non-zero",
			input:     entity.Task{Type: "telegram", Points: 100},
			expect:    entity.Task{},
			expectErr: entity.ErrInvalidTaskPoints,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := tt.input.Validate()
			if tt.expectErr != nil {
				require.ErrorIs(t, err, tt.expectErr)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expect, task)
		})
	}
}
