package todo

import (
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
	Update(id Id, title, description string) (Todo, error)
	Delete(id Id) error
	Count() (uint, error)
}
