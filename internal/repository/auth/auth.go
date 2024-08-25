package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"cmd/app/main.go/internal/model"
	repo "cmd/app/main.go/internal/repository"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, u model.User) error {
	query := ` 
		INSERT INTO users (id, nickname, email, password, role)
 		SELECT $1, LOWER($2), LOWER($3), $4, LOWER($5)
		RETURNING id, nickname, email, password, role`

	_, err := r.pool.Exec(ctx, query, u.ID, u.Nickname, u.Email, u.Password, u.Role)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return repo.ErrDuplicateField{Field: pgErr.ConstraintName}
			}
		}

		return errors.WithStack(err)
	}

	return nil
}
