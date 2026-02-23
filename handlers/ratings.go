package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/slikasp/fragrancetrackgo/internal/database/localDatabase"
)

// create a new user-owned rating entry. Flow:
// 1. Normalize brand/name (trim + lowercase) so lookups are consistent.
// 2. Validate required keys.
// 3. Check if an entry already exists for (user, brand, name).
// 4. Insert the new record when no duplicate is found.
func RatingAdd(ctx context.Context, store RatingStore, userID uuid.UUID, input RatingInput) (localDatabase.Rating, error) {
	brand := strings.ToLower(strings.TrimSpace(input.Brand))
	name := strings.ToLower(strings.TrimSpace(input.Name))
	if brand == "" || name == "" {
		return localDatabase.Rating{}, ErrInvalidInput
	}

	_, err := store.GetRating(ctx, localDatabase.GetRatingParams{
		UserID: userID,
		Brand:  brand,
		Name:   name,
	})
	if err == nil {
		return localDatabase.Rating{}, ErrRatingExists
	}
	// Any lookup error besides "not found" is an operational failure.
	if !errors.Is(err, sql.ErrNoRows) {
		return localDatabase.Rating{}, fmt.Errorf("lookup rating: %w", err)
	}

	newFrag, err := store.AddRating(ctx, localDatabase.AddRatingParams{
		UserID:  userID,
		Brand:   brand,
		Name:    name,
		Rating:  input.Rating,
		Comment: input.Comment,
	})
	if err != nil {
		return localDatabase.Rating{}, fmt.Errorf("add rating: %w", err)
	}

	return newFrag, nil
}

// delete one rating identified by (user, brand, name) and return it
func RatingRemove(ctx context.Context, store RatingStore, userID uuid.UUID, brand, name string) (localDatabase.Rating, error) {
	removedFrag, err := store.RemoveRating(ctx, localDatabase.RemoveRatingParams{
		UserID: userID,
		Brand:  brand,
		Name:   name,
	})
	if err != nil {
		return localDatabase.Rating{}, fmt.Errorf("remove rating: %w", err)
	}

	return removedFrag, nil
}

// update optional fields (comment/rating) on one rating row
// sql.Null* types preserve "unset" semantics through to SQL
func RatingUpdate(ctx context.Context, store RatingStore, userID uuid.UUID, brand, name string, comment sql.NullString, rating sql.NullInt32) (localDatabase.Rating, error) {
	updatedFrag, err := store.UpdateRating(ctx, localDatabase.UpdateRatingParams{
		UserID:  userID,
		Brand:   brand,
		Name:    name,
		Rating:  rating,
		Comment: comment,
	})
	if err != nil {
		return localDatabase.Rating{}, fmt.Errorf("update rating: %w", err)
	}

	return updatedFrag, nil
}

// RatingList returns all ratings for a specific user
func RatingList(ctx context.Context, store RatingStore, userID uuid.UUID) ([]localDatabase.Rating, error) {
	frags, err := store.GetRatings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list ratings: %w", err)
	}

	return frags, nil
}
