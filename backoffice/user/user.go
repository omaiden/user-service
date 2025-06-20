package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"user-service/backoffice/internal/kctx"
	"user-service/pkg/sql"

	dbsql "database/sql"
	"github.com/acoshift/pgsql/pgctx"
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

type Article struct {
	ID        int
	Title     string
	Content   string
	AuthorID  string
	CreatedAt time.Time
	UpdatedAt time.Time
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

func GetUserArticles(ctx context.Context, userID string, limit int) ([]*Article, error) {
	if limit <= 0 || limit > 100 {
		return nil, errors.New("limit must be between 1 and 100")
	}

	currentUserID := kctx.GetUserID(ctx)
	if currentUserID == "" {
		return nil, errors.New("unauthorized")
	}

	rows, err := pgctx.Query(ctx, `
        SELECT id, title, content, author_id, created_at, updated_at
        FROM articles 
        WHERE author_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []*Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.AuthorID, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, &a)
	}

	return articles, nil
}
