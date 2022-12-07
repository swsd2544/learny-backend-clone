package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/swsd2544/learny-backend-clone/internal/entity"
	"log"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (r UserRepository) Insert(user *entity.User) error {
	query := `INSERT INTO users (username, firstname, lastname, email, hash_password,
    coin, users.role) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, 
    created_at, version`

	args := []any{user.Username, user.Firstname, user.Lastname, user.Email,
		user.HashPassword, user.Coin, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		log.Println(err)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r UserRepository) GetUsersWithClassID(classID int64) ([]*entity.User, error) {
	query := `SELECT users.id, users.firstname, users.lastname, users.email, 
    users.hash_password, users.coin, users.role, users.version FROM users INNER JOIN enrollments
    ON users.id = enrolls.user_id WHERE enrolls.class_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := r.db.Query(ctx, query, classID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return []*entity.User{}, nil
		default:
			return nil, err
		}
	}
	defer results.Close()

	var users []*entity.User
	for results.Next() {
		var user entity.User
		err := results.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = results.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r UserRepository) GetUsersWithID(userID int64) (*entity.User, error) {
	query := `SELECT id, firstname, lastname, email, hash_password, 
       coin, users.role, version FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user entity.User
	err := r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Firstname, &user.Lastname,
		&user.Email, &user.HashPassword, &user.Coin, &user.Role, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r UserRepository) GetUsersWithEmail(email string) (*entity.User, error) {
	query := `SELECT id, username, firstname, lastname, hash_password, 
       coin, users.role, version, character_id FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Firstname, &user.Lastname,
		&user.HashPassword, &user.Coin, &user.Role, &user.Version, &user.CharacterID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r UserRepository) Update(user *entity.User) error {
	query := `UPDATE users SET firstname=$1, lastname=$2, email=$3, hash_password=$4, 
    coin=$5, users.role=$6, version = version + 1 WHERE id = $7 AND version = $8 RETURNING version`

	args := []any{
		user.Firstname,
		user.Lastname,
		user.Email,
		user.HashPassword,
		user.Coin,
		user.Role,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRow(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}