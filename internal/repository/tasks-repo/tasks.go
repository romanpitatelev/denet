package tasksrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/romanpitatelev/denet/internal/entity"
	"github.com/romanpitatelev/denet/internal/repository/store"
)

type Repo struct {
	db *store.DataStore
}

func New(db *store.DataStore) *Repo {
	return &Repo{
		db: db,
	}
}

const (
	insertTaskQuery = `
INSERT INTO tasks (id, user_id, type, created_at)
VALUES ($1, $2, $3, $4)
`

	updateUserPointsQuery = `
UPDATE users
SET points = points + $1,
	updated_at = $2
WHERE id = $3
RETURNING points
`
)

func (r *Repo) Task(ctx context.Context, userID entity.UserID, task entity.Task) (entity.TaskResponse, error) {
	var response entity.TaskResponse
	transactionTime := time.Now()

	if err := r.db.WithinTransaction(ctx, func(ctx context.Context, tx store.Transaction) error {
		_, err := tx.Exec(ctx, insertTaskQuery,
			task.ID,
			userID,
			task.Type,
			transactionTime,
		)
		if err != nil {
			return fmt.Errorf("failed to insert task: %w", err)
		}

		row := tx.QueryRow(ctx, updateUserPointsQuery,
			task.Points,
			transactionTime,
			userID,
		)

		var totalPoints int

		err = row.Scan(
			&totalPoints,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return entity.ErrUserNotFound
			}
			return fmt.Errorf("failed to update user points: %w", err)
		}

		response = entity.TaskResponse{
			Task: entity.Task{
				ID:     task.ID,
				Type:   task.Type,
				Points: task.Points,
			},
			CreatedAt:    transactionTime,
			TotatlPoints: totalPoints,
		}

		return nil
	}); err != nil {
		return entity.TaskResponse{}, err
	}

	return response, nil
}
