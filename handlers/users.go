package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/slikasp/fragrancetrackgo/internal/auth"
	"github.com/slikasp/fragrancetrackgo/internal/database/localDatabase"
)

func UserRegister(ctx context.Context, users UserStore, name, pass string) (localDatabase.User, error) {
	name = strings.TrimSpace(name)
	if name == "" || pass == "" {
		return localDatabase.User{}, ErrInvalidInput
	}

	_, err := users.GetUserByName(ctx, name)
	if err == nil {
		return localDatabase.User{}, ErrUserAlreadyExists
	}
	// Any lookup error besides "not found" is treated as infrastructure failure.
	if !errors.Is(err, sql.ErrNoRows) {
		return localDatabase.User{}, fmt.Errorf("lookup user by name: %w", err)
	}

	// Store only hashed passwords
	hpass, err := auth.HashPassword(pass)
	if err != nil {
		return localDatabase.User{}, fmt.Errorf("hash password: %w", err)
	}

	newUser, err := users.CreateUser(ctx, localDatabase.CreateUserParams{
		Name:           name,
		HashedPassword: hpass,
	})
	if err != nil {
		return localDatabase.User{}, fmt.Errorf("create user %q: %w", name, err)
	}

	return newUser, nil
}

func UserLogin(ctx context.Context, users UserStore, name, pass string) (localDatabase.User, error) {
	name = strings.TrimSpace(name)
	if name == "" || pass == "" {
		return localDatabase.User{}, ErrInvalidInput
	}

	user, err := users.GetUserByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return localDatabase.User{}, ErrInvalidCredential
		}
		return localDatabase.User{}, fmt.Errorf("lookup user by name: %w", err)
	}

	// Password hash comparison
	match, err := auth.CheckPasswordHash(pass, user.HashedPassword)
	if err != nil {
		return localDatabase.User{}, fmt.Errorf("check password hash: %w", err)
	}
	if !match {
		return localDatabase.User{}, ErrInvalidCredential
	}

	return user, nil
}

// TODO: find a use for it (maybe admin console?) or get rid of it
func UserList(ctx context.Context, users UserStore, currentUserID uuid.UUID) ([]string, error) {
	names, err := users.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return names, nil
}
