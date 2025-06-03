package entity

import "errors"

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidTaskPoints    = errors.New("user cannot set the points")
	ErrInvalidTaskType      = errors.New("task type has be either 'telegram' or 'twitter'")
	ErrInvalidUUIDFormat    = errors.New("invalid UUID format")
	ErrDuplicateContact     = errors.New("duplicate contact")
	ErrSelfReference        = errors.New("self reference prohibited")
	ErrEmptyReferenceUser   = errors.New("error empty reference user ID")
)
