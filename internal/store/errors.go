package store

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrUserAlreadyExists = errors.New("user with this username or email already exists")
)
