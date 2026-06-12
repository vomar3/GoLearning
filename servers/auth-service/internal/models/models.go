package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserBadResponse struct {
	Error string `json:"error"`
	Msg   string `json:"message"`
}

type UserGoodResponse struct {
	UserID int `json:"userID"`
}

type Authorization struct {
	UserID       int    `json:"userID"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

func (user *User) Validate() bool {
	if user.Email == "" || user.Password == "" {
		return false
	}

	return true
}
