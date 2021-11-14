package service

type User struct {
}

type ServiceInterface interface {
	Login(userId string, password string) error
	ListUsers() ([]User, error)
	Logout(userId string)
}
