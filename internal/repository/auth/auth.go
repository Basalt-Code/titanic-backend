package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"cmd/app/main.go/internal/model/domain"
)

var refreshExpirationTime = time.Now().AddDate(0, 1, 0).Format(time.RFC3339)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) CreateSession(ctx context.Context, session domain.Session) error {
	query := `
		INSERT INTO user_sessions (user_id, session_id, refresh_token, expired_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(ctx, query, session.UserID, session.ID, session.RefreshToken, refreshExpirationTime)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) UpdateRefreshToken(ctx context.Context, session domain.Session) error {
	query := `
		UPDATE user_sessions 
		SET refresh_token = $1, 
		    expired_at = $2
		WHERE user_id = $3
			AND session_id = $4
			AND deleted_at IS NULL
			AND expired_at > CURRENT_DATE
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		session.RefreshToken,
		refreshExpirationTime,
		session.UserID,
		session.ID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) GetRefreshToken(ctx context.Context, sessionID string) (*domain.RefreshTokens, error) {
	query := `
		SELECT refresh_token, expired_at
		FROM user_sessions
		WHERE session_id = $1
		AND deleted_at IS NULL
		ANd expired_at > CURRENT_DATE
	`

	var refreshToken domain.RefreshTokens
	if err := r.pool.QueryRow(ctx, query, sessionID).Scan(
		&refreshToken.Refresh,
		&refreshToken.ExpiredAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.WithStack(err)
	}

	return &refreshToken, nil
}

func (r *Repo) DeleteRefreshToken(ctx context.Context, userID, sessionID string) error {
	query := `
		UPDATE user_sessions 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE user_id = $1 AND session_id = $2
	`

	if _, err := r.pool.Exec(ctx, query, userID, sessionID); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
