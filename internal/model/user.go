package model

type RegistrationCredentials struct {
	Nickname string
	Email    string
	Password string
}

type User struct {
	ID       string
	Nickname *string
	Email    *string
	Password *string
	Role     *string
}
