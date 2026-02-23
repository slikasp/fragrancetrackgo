package handlers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/slikasp/fragrancetrackgo/internal/database/localDatabase"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrRatingExists      = errors.New("rating already exists")
)

type UserStore interface {
	CreateUser(context.Context, localDatabase.CreateUserParams) (localDatabase.User, error)
	GetUserByName(context.Context, string) (localDatabase.User, error)
	GetUserByID(context.Context, uuid.UUID) (localDatabase.User, error)
	GetUsers(context.Context) ([]string, error)
	ResetUsers(context.Context) error
}

type RatingStore interface {
	GetRating(context.Context, localDatabase.GetRatingParams) (localDatabase.Rating, error)
	AddRating(context.Context, localDatabase.AddRatingParams) (localDatabase.Rating, error)
	RemoveRating(context.Context, localDatabase.RemoveRatingParams) (localDatabase.Rating, error)
	UpdateRating(context.Context, localDatabase.UpdateRatingParams) (localDatabase.Rating, error)
	GetRatings(context.Context, uuid.UUID) ([]localDatabase.Rating, error)
}

type RatingInput struct {
	Brand   string
	Name    string
	Rating  sql.NullInt32
	Comment sql.NullString
}
