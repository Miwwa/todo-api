package user

type Id string

type User struct {
	Id           Id     `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	passwordHash string
}
