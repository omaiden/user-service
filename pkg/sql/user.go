package sql

import (
	"context"
	"github.com/acoshift/pgsql/pgctx"
)

type User struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"unique"`
	Email string `gorm:""`
}

func (u *User) TableName() string {
	return "user"
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := pgctx.QueryRow(ctx, `
	select *
		from "user"
		where email = $1
		`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	return &user, err
}

func CreateUser(ctx context.Context, user *User) error {
	return pgctx.QueryRow(ctx, `
        INSERT INTO "user" (name, email) VALUES ($1, $2)
        RETURNING id
    `, user.Name, user.Email).Scan(&user.ID)
}
