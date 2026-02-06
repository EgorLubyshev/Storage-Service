package postgres

import (
	"context"
	"database/sql"

	"github.com/your-username/storage-service/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, name, passwordHash string) (string, error) {
	var id string
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (name, password_hash) VALUES ($1, $2) RETURNING id`,
		name,
		passwordHash,
	).Scan(&id)
	return id, err
}

func (r *UserRepo) GetByName(ctx context.Context, name string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, password_hash, rights, created_at FROM users WHERE name = $1`,
		name,
	).Scan(&user.ID, &user.Name, &user.PasswordHash, &user.Rights, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
