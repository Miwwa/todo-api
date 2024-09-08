package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"strings"
)

var (
	AlreadyExists = errors.New("user already exists")
	NotFound      = errors.New("user not found")
)

type Storage interface {
	Create(ctx context.Context, email string, passwordHash string, name string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type SqliteUsersStorage struct {
	db *sql.DB
}

func NewSqliteUsersStorage(db *sql.DB) *SqliteUsersStorage {
	return &SqliteUsersStorage{db: db}
}

func (s SqliteUsersStorage) Create(ctx context.Context, email string, passwordHash string, name string) (User, error) {
	userId := ulid.Make().String()

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (id, email, name, password_hash) VALUES (?, ?, ?, ?)")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	exec, err := stmt.ExecContext(ctx, userId, email, name, passwordHash)
	if err != nil {
		return User{}, mapError(err)
	}

	affected, err := exec.RowsAffected()
	if err != nil {
		return User{}, err
	}
	if affected == 0 {
		return User{}, errors.New("user creation error")
	}

	return User{
		Id:    Id(userId),
		Email: email,
		Name:  name,
	}, nil
}

func (s SqliteUsersStorage) GetUserByEmail(ctx context.Context, email string) (User, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT id, email, name, password_hash FROM users WHERE email=?")
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	var user = User{}
	err = stmt.QueryRowContext(ctx, email).Scan(&user.Id, &user.Email, &user.Name, &user.passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, NotFound
		}
		return User{}, err
	}

	return user, nil
}

func mapError(err error) error {
	var sqlErr sqlite3.Error
	if errors.As(err, &sqlErr) {
		if errors.Is(sqlErr.Code, sqlite3.ErrConstraint) && strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return AlreadyExists
		}
	}
	return err
}
