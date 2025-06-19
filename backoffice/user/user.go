package user

import (
	"context"
	"errors"
	"strings"

	"user-service/pkg/sql"

	dbsql "database/sql"
	"github.com/moonrhythm/validator"
)

type CreateUserRequest struct {
	Name  string
	Email string
}

type User struct {
	ID    uint64
	Name  string
	Email string
}

func (p *CreateUserRequest) Valid() error {
	// normalize data
	p.Name = strings.TrimSpace(p.Name)
	p.Email = strings.TrimSpace(p.Email)

	// validate data
	v := validator.New()
	v.Must(p.Name != "", "name required")
	v.Must(p.Email != "", "email required")
	return v.Error()
}

func CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// Validate request
	if err := req.Valid(); err != nil {
		return nil, err
	}

	// Check if email exists
	_, err := sql.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(dbsql.ErrNoRows, err) {
		return nil, err
	}
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Save user
	sqlUser := sql.User{
		Name:  req.Name,
		Email: req.Email,
	}
	err = sql.CreateUser(ctx, &sqlUser)
	if err != nil {
		return nil, err
	}

	// Send welcome email
	//if err := h.emailSvc.SendWelcomeEmail(ctx, user); err != nil {
	//	// Don't fail the entire operation for email
	//}

	return &User{
		ID:    sqlUser.ID,
		Name:  sqlUser.Name,
		Email: sqlUser.Email,
	}, nil
}
