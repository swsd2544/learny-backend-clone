package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swsd2544/learny-backend-clone/internal/entity"
	"time"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func (r TokenRepository) New(userID int64, ttl time.Duration, scope string) (*entity.Token, error) {
	token, err := entity.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = r.Insert(token)
	return token, err
}

func (r TokenRepository) Insert(token *entity.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r TokenRepository) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens WHERE scope=$1 AND user_id=$2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query, scope, userID)
	return err
}
