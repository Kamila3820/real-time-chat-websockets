package user

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int

	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	// fmt.Printf("Executing query: %s\nParams: username=%s, password=%s, email=%s\n", query, user.Username, user.Password, user.Email)

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&lastInsertId)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	user.ID = int64(lastInsertId)

	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}

	query := "SELECT id, email, username, password FROM users WHERE email = $1"

	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return &User{}, nil
	}

	return &user, nil
}
