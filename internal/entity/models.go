package entity

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	UserID      uuid.UUID
	TaskID      uuid.UUID
	ReferenceID uuid.UUID
)

type User struct {
	ID        UserID    `json:"userId"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Claims struct {
	UserID    UserID     `json:"userId"`
	Name      string     `json:"name"`
	Email     *string    `json:"email"`
	Role      *string    `json:"role"`
	Points    int        `json:"points"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	jwt.RegisteredClaims
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

func (uu *UserUpdate) Validate() (UserUpdate, error) {
	if uu.Email != nil {
		if err := validateEmail(*uu.Email); err != nil {
			return UserUpdate{}, fmt.Errorf("invalid email address: %w", err)
		}
	}

	return *uu, nil
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("email validation error: %w", err)
	}

	return nil
}

type ListRequest struct {
	Sorting    string
	Descending bool
	Limit      int
	Filter     string
	Offset     int
}

type Task struct {
	ID     TaskID `json:"id"`
	Type   string `json:"type"`
	Points int    `json:"points"`
}

type TaskResponse struct {
	Task
	CreatedAt    time.Time `json:"createdAt"`
	TotatlPoints int       `json:"totalPoints"`
}

func (t *Task) Validate() (Task, error) {
	if t.Points != 0 {
		return Task{}, ErrInvalidTaskPoints
	}

	if t.Type == "telegram" {
		t.Points = 3

		return *t, nil
	}

	if t.Type == "twitter" {
		t.Points = 2

		return *t, nil
	}

	return Task{}, ErrInvalidTaskType
}

type Reference struct {
	ID              ReferenceID `json:"referenceId"`
	UserID          UserID      `json:"userId"`
	UserReferenceID UserID      `json:"userReferenceId"`
	Points          int         `json:"points"`
}

type ReferenceResponse struct {
	Reference
	CreatedAt    time.Time `json:"createdAt"`
	TotatlPoints int       `json:"totalPoints"`
}

func (r *Reference) Validate() error {
	if r.UserID == r.UserReferenceID {
		return ErrSelfReference
	}

	if r.UserReferenceID == UserID(uuid.Nil) {
		return ErrEmptyReferenceUser
	}

	return nil
}

func unmarshalUUID(id *uuid.UUID, data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("unmarshalling error: %w", err)
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return ErrInvalidUUIDFormat
	}

	*id = parsed

	return nil
}

func (u *UserID) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), data)
}

func (u *TaskID) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), data)
}

func (u *ReferenceID) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), data)
}

//nolint:wrapcheck
func (u UserID) MarshalText() ([]byte, error) {
	return json.Marshal(uuid.UUID(u).String())
}

//nolint:wrapcheck
func (u TaskID) MarshalText() ([]byte, error) {
	return json.Marshal(uuid.UUID(u).String())
}

//nolint:wrapcheck
func (u ReferenceID) MarshalText() ([]byte, error) {
	return json.Marshal(uuid.UUID(u).String())
}
