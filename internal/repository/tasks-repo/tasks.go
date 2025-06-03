package tasksrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	insertReferenceQuery = `
INSERT INTO reference (id, user_id, reference_id, created_at)
VALUES ($1, $2, $3, $4)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return entity.TaskResponse{}, entity.ErrUserNotFound
		}

		return entity.TaskResponse{}, fmt.Errorf("failed to complete task transaction: %w", err)
	}

	return response, nil
}

func (r *Repo) ReferralTask(ctx context.Context, reference entity.Reference) (entity.ReferenceResponse, error) {
	var response entity.ReferenceResponse

	transactionTime := time.Now()

	if err := r.db.WithinTransaction(ctx, func(ctx context.Context, tx store.Transaction) error {
		_, err := tx.Exec(ctx, insertReferenceQuery,
			reference.ID,
			reference.UserID,
			reference.UserReferenceID,
			transactionTime,
		)
		if err != nil {
			return fmt.Errorf("failed to insert task: %w", err)
		}

		row := tx.QueryRow(ctx, updateUserPointsQuery,
			reference.Points,
			transactionTime,
			reference.UserID,
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

		response = entity.ReferenceResponse{
			Reference: entity.Reference{
				ID:              reference.ID,
				UserID:          reference.UserID,
				UserReferenceID: reference.UserReferenceID,
				Points:          reference.Points,
			},
			CreatedAt:    transactionTime,
			TotatlPoints: totalPoints,
		}

		return nil
	}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return entity.ReferenceResponse{}, entity.ErrUserNotFound
		}

		return entity.ReferenceResponse{}, fmt.Errorf("failed to complete reference transaction: %w", err)
	}

	return response, nil
}
