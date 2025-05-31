package usersrepo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/romanpitatelev/denet/internal/entity"
	"github.com/romanpitatelev/denet/internal/repository/store"
)

const (
	maxUpdates = 3
)

type Repo struct {
	db *store.DataStore
}

func New(db *store.DataStore) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetUser(ctx context.Context, userID entity.UserID) (entity.User, error) {
	db := r.db.GetTXFromContext(ctx)

	var user entity.User

	query := `
SELECT id, name, email, role, points, created_at, updated_at	
FROM users
WHERE TRUE
	AND id = $1
	AND deleted_at IS NULL`

	if err := pgxscan.Get(ctx, db, &user, query, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}

		return entity.User{}, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

func (r *Repo) UpdateUser(ctx context.Context, userID entity.UserID, updatedUser entity.UserUpdate) (entity.User, error) {
	db := r.db.GetTXFromContext(ctx)

	var (
		sb     strings.Builder
		params []interface{}
	)

	updates := make([]string, 0, maxUpdates)

	sb.WriteString("UPDATE users SET ")

	if updatedUser.Name != nil {
		params = append(params, *updatedUser.Name)
		updates = append(updates, fmt.Sprintf("name = $%d", len(params)))
	}

	if updatedUser.Email != nil {
		params = append(params, *updatedUser.Email)
		updates = append(updates, fmt.Sprintf("email = $%d", len(params)))
	}

	if updatedUser.Role != nil {
		params = append(params, *updatedUser.Role)
		updates = append(updates, fmt.Sprintf("role = $%d", len(params)))
	}

	if len(updates) == 0 {
		return r.GetUser(ctx, userID)
	}

	params = append(params, time.Now())
	updates = append(updates, fmt.Sprintf("updated_at = $%d", len(params)))

	sb.WriteString(strings.Join(updates, ", "))

	params = append(params, userID)

	sb.WriteString(" RETURNING id, name, email, role, points, created_at, updated_at")

	var user entity.User

	if err := pgxscan.Get(ctx, db, &user, sb.String(), params...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, entity.ErrUserNotFound
		}
	}

	return user, nil
}

func (r *Repo) DeleteUser(ctx context.Context, userID entity.UserID) error {
	_, err := r.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return entity.ErrUserNotFound
		}

		return fmt.Errorf("failed to fetch user in DeleteUser(): %w", err)
	}

	query := `
UPDATE users
SET deleted_at = NOW()
WHERE TRUE
	AND id = $1
	AND deleted_at IS NULL`

	_, err = r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user %s: %w", uuid.UUID(userID), err)
	}

	return nil
}

const defaultLimit = 25

const listUsersQuery = `
SELECT id, name, email, role, points, created_at, updated_at
FROM users
WHERE deleted_at IS NULL
`

func (r *Repo) GetUsers(ctx context.Context, request entity.ListRequest) ([]entity.User, error) {
	db := r.db.GetTXFromContext(ctx)

	mapping := map[string]string{
		"name":   "name",
		"points": "points",
	}

	var args []any

	sb := strings.Builder{}
	sb.WriteString(listUsersQuery)

	if request.Filter != "" {
		args = append(args, "%"+request.Filter+"%")
		sb.WriteString(fmt.Sprintf(` AND concat_ws('', name, role, points, created_at, updated_at) ILIKE $%d`, len(args)))
	}

	orderBy := mapping[request.Sorting]
	if orderBy == "" {
		orderBy = mapping["name"]
	}

	sb.WriteString(" ORDER BY " + orderBy)

	if request.Descending {
		sb.WriteString(" DESC")
	}

	limit := defaultLimit
	if request.Limit > 0 {
		limit = request.Limit
	}

	sb.WriteString(" LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(request.Offset))

	var result []entity.User

	if err := pgxscan.Select(ctx, db, &result, sb.String(), args...); err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return result, nil
}
