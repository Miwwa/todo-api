package todo

import (
	"database/sql"
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
	Create(userId user.Id, title, description string) (Todo, error)
	Get(userId user.Id, limit, offset uint) ([]Todo, error)
	Update(userId user.Id, id Id, title, description string) (Todo, error)
	Delete(userId user.Id, id Id) error
	Count() (uint, error)
}

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(db *sql.DB) *SqliteStorage {
	return &SqliteStorage{db: db}
}

func (s SqliteStorage) Create(userId user.Id, title, description string) (Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (s SqliteStorage) Get(userId user.Id, limit, offset uint) ([]Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (s SqliteStorage) Update(userId user.Id, id Id, title, description string) (Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (s SqliteStorage) Delete(userId user.Id, id Id) error {
	//TODO implement me
	panic("implement me")
}

func (s SqliteStorage) Count() (uint, error) {
	//TODO implement me
	panic("implement me")
}
