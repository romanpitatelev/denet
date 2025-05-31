package entity

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type (
	UserID uuid.UUID
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
