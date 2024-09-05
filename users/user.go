package users

type User struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	passwordHash string
}

type RegistrationData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
