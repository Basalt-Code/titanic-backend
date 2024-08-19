package store

import "cmd/app/main.go/internal/model"

type UserRepository interface {
	Create(*model.User) error
	FindByUsername(string) (*model.User, error)
	Find(int) (*model.User, error)
}
