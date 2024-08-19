package sqlstore

import (
	"cmd/app/main.go/internal/store"
	"database/sql"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(storagePath string) (*Store, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{store: s}
	return s.userRepository
}
