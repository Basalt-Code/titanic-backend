package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"cmd/app/main.go/internal/model/domain"
	"cmd/app/main.go/internal/model/dto"
	repo "cmd/app/main.go/internal/repository"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, u dto.User) error {
	query := `
		INSERT INTO users (id, username, email, password, role)
		VALUES ($1, LOWER($2), LOWER($3), $4, LOWER($5))`

	_, err := r.pool.Exec(ctx, query, u.ID, u.Username, u.Email, u.Password, u.Role)
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

func (r *Repo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password, role
		FROM users
		WHERE deleted_at IS NULL AND username = $1`

	var user domain.User
	if err := r.pool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.WithStack(err)
	}

	return &user, nil
}
