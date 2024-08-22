package sqlstore

import (
	"cmd/app/main.go/internal/model"
	"cmd/app/main.go/internal/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	user, _ := r.FindByEmail(u.Email)
	if user != nil {
		return store.ErrUserAlreadyExists
	}
	user, _ = r.FindByUsername(u.Username)
	if user != nil {
		return store.ErrUserAlreadyExists
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (username, encrypted_password, email, role) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Username,
		u.EncryptedPassword,
		u.Email,
		u.Role,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, username, email, encrypted_password, role FROM users WHERE username = $1",
		username,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.EncryptedPassword,
		&u.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, username, email, encrypted_password, role FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.EncryptedPassword,
		&u.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, username, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err

	}
	return u, nil
}
