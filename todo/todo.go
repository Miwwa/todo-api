package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"time"
	"todo-api/user"
)

type Id string

type Todo struct {
	Id          Id        `json:"id"`
	UserId      user.Id   `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Todo) Invalid() bool {
	return t.Id == ""
}

type Storage interface {
	Create(ctx context.Context, userId user.Id, title, description string) (Todo, error)
	GetById(ctx context.Context, id Id) (Todo, error)
	GetByUserId(ctx context.Context, userId user.Id, options FindOptions) ([]Todo, error)
	Update(ctx context.Context, id Id, title, description string) (Todo, error)
	Delete(ctx context.Context, id Id) error
	Count(ctx context.Context, userId user.Id) (uint, error)
}

const (
	IdName          = "id"
	TitleName       = "title"
	DescriptionName = "description"

	SortAscending  = "asc"
	SortDescending = "desc"
)

type FindOptions struct {
	Limit, Offset      uint
	SortBy, SortOrder  string
	Title, Description string
}

func (f *FindOptions) Validate() error {
	if f.Limit < 0 {
		return errors.New("limit cannot be negative")
	}
	if f.Offset < 0 {
		return errors.New("offset cannot be negative")
	}
	if f.SortOrder != SortAscending && f.SortOrder != SortDescending {
		return errors.New("invalid sort order")
	}
	if f.SortBy != IdName && f.SortBy != TitleName && f.SortBy != DescriptionName {
		return errors.New("invalid sort field")
	}
	return nil
}

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(db *sql.DB) *SqliteStorage {
	return &SqliteStorage{db: db}
}

func (s SqliteStorage) Create(ctx context.Context, userId user.Id, title, description string) (Todo, error) {
	todoId := ulid.Make().String()

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO todos (id, user_id, title, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return Todo{}, err
	}
	defer stmt.Close()

	exec, err := stmt.ExecContext(ctx, todoId, userId, title, description)
	if err != nil {
		return Todo{}, err
	}

	affected, err := exec.RowsAffected()
	if affected == 0 || err != nil {
		return Todo{}, errors.New("user creation error")
	}

	return Todo{
		Id:          Id(todoId),
		UserId:      userId,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s SqliteStorage) GetById(ctx context.Context, id Id) (Todo, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT id, user_id, title, description, created_at, updated_at
		FROM todos WHERE id=?
	`)
	if err != nil {
		return Todo{}, err
	}
	defer stmt.Close()
	todo := Todo{}

	err = stmt.QueryRowContext(ctx, id).Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, nil
		}
		return Todo{}, err
	}
	return todo, nil
}

func (s SqliteStorage) GetByUserId(ctx context.Context, userId user.Id, options FindOptions) ([]Todo, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}

	whereTitle := ""
	whereDescription := ""
	if options.Title != "" {
		whereTitle += " AND title LIKE ?"
	}
	if options.Description != "" {
		whereDescription += " AND description LIKE ?"
	}

	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf(`
		SELECT id, user_id, title, description, created_at, updated_at
		FROM todos WHERE user_id=? %s %s
		ORDER BY %s %s
		LIMIT ?
		OFFSET ?
	`, whereTitle, whereDescription, options.SortBy, options.SortOrder))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := []any{userId}
	if options.Title != "" {
		args = append(args, options.Title)
	}
	if options.Description != "" {
		args = append(args, options.Description)
	}
	args = append(args, options.Limit, options.Offset)

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s SqliteStorage) Update(ctx context.Context, id Id, title, description string) (Todo, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		UPDATE todos
		SET title=?,description=?,updated_at=CURRENT_TIMESTAMP
		WHERE id=?
		RETURNING id, user_id, title, description, created_at, updated_at
	`)
	if err != nil {
		return Todo{}, err
	}
	defer stmt.Close()

	todo := Todo{}
	err = stmt.QueryRowContext(ctx, title, description, id).
		Scan(&todo.Id, &todo.UserId, &todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Todo{}, nil
		}
		return Todo{}, err
	}

	return todo, nil
}

func (s SqliteStorage) Delete(ctx context.Context, id Id) error {
	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM todos WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s SqliteStorage) Count(ctx context.Context, userId user.Id) (uint, error) {
	stmt, err := s.db.PrepareContext(ctx, "SELECT COUNT() FROM todos WHERE user_id=?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count uint
	err = stmt.QueryRowContext(ctx, userId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
