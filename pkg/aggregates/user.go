package aggregates

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(email string, name string, password string) (User, error) {
	user := User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	return user, nil
}
