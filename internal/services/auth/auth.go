package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"

	"cmd/app/main.go/internal/config"
	"cmd/app/main.go/internal/model/domain"
	"cmd/app/main.go/internal/model/dto"
)

type userRepo interface {
	Create(ctx context.Context, u dto.User) error
	FindByUsername(ctx context.Context, email string) (*domain.User, error)
}

type authRepo interface {
	CreateSession(ctx context.Context, session domain.Session) error
	GetRefreshToken(ctx context.Context, sessionID string) (*domain.RefreshTokens, error)
	UpdateRefreshToken(ctx context.Context, session domain.Session) error
	DeleteRefreshToken(ctx context.Context, userID, sessionID string) error
}

type Service struct {
	cfg      config.ServerConfig
	userRepo userRepo
	authRepo authRepo
}

var (
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrInvalidTokenPair    = errors.New("invalid token pair")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

func New(cfg config.ServerConfig, userRepo userRepo, authRepo authRepo) *Service {
	return &Service{
		cfg:      cfg,
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *Service) Register(ctx context.Context, credentials dto.RegistrationCredentials) error {
	passHash, err := s.hash(credentials.Password)
	if err != nil {
		return err
	}

	err = s.userRepo.Create(ctx, dto.User{
		ID:       uuid.NewString(),
		Username: &credentials.Username,
		Email:    &credentials.Email,
		Password: lo.ToPtr(string(passHash)),
		Role:     &credentials.Role,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, credentials dto.Credentials) (domain.Tokens, error) {
	user, err := s.userRepo.FindByUsername(ctx, credentials.Username)
	if err != nil {
		return domain.Tokens{}, err
	}

	if user == nil {
		return domain.Tokens{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return domain.Tokens{}, ErrInvalidCredentials
	}

	return s.createTokens(ctx, *user)
}

func (s *Service) Logout(ctx context.Context, userID, sessionID string) error {
	return s.authRepo.DeleteRefreshToken(ctx, userID, sessionID)
}

func (s *Service) RefreshToken(ctx context.Context, oldTokens domain.Tokens) (domain.Tokens, error) {
	if !isValidTokenPair(oldTokens) {
		return domain.Tokens{}, ErrInvalidTokenPair
	}

	claims, err := parseUnverifiedToken(oldTokens.Access)
	if err != nil {
		return domain.Tokens{}, err
	}

	err = s.verifyRefreshToken(ctx, claims.SessionID, oldTokens.Refresh)
	if err != nil {
		return domain.Tokens{}, err
	}

	tokens, err := s.newTokens(domain.CustomClaims{
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
		Role:      claims.Role,
	})
	if err != nil {
		return domain.Tokens{}, err
	}

	err = s.authRepo.UpdateRefreshToken(ctx, domain.Session{
		ID:           claims.SessionID,
		UserID:       claims.UserID,
		RefreshToken: tokens.Refresh,
	})
	if err != nil {
		return domain.Tokens{}, err
	}

	return tokens, nil
}

func (s *Service) ParseTokenWitClaims(access string) (*jwt.Token, *domain.Claims, error) {
	token, err := jwt.ParseWithClaims(access, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.SecretKey), nil
	})

	if token != nil {
		if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
			return token, claims, nil
		}
	}

	return nil, nil, err
}

func (s *Service) createTokens(ctx context.Context, user domain.User) (domain.Tokens, error) {
	sessionID := uuid.New().String()
	tokens, err := s.newTokens(domain.CustomClaims{
		UserID:    user.ID,
		SessionID: sessionID,
		Role:      user.Role,
	})
	if err != nil {
		return domain.Tokens{}, err
	}

	err = s.authRepo.CreateSession(ctx, domain.Session{
		ID:           sessionID,
		UserID:       user.ID,
		RefreshToken: tokens.Refresh,
	})
	if err != nil {
		return domain.Tokens{}, err
	}

	return tokens, nil
}

func (s *Service) newTokens(claims domain.CustomClaims) (domain.Tokens, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		CustomClaims: claims,
	})

	accessToken, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return domain.Tokens{}, err
	}

	return domain.Tokens{
		Access:  accessToken,
		Refresh: newRefreshToken(accessToken),
	}, nil
}

func newRefreshToken(accessToken string) string {
	return fmt.Sprintf(uuid.New().String() + accessToken[len(accessToken)-6:])
}

func (s *Service) hash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}

func (s *Service) verifyRefreshToken(ctx context.Context, sessionId, token string) error {
	refresh, err := s.authRepo.GetRefreshToken(ctx, sessionId)
	if err != nil {
		return err
	}

	if refresh == nil {
		return ErrInvalidRefreshToken
	}

	if token != refresh.Refresh {
		return ErrInvalidRefreshToken
	}

	return nil
}

func parseUnverifiedToken(accessToken string) (*domain.Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, &domain.Claims{})
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, nil
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	return claims, nil
}

func isValidTokenPair(tokens domain.Tokens) bool {
	return tokens.Access[len(tokens.Access)-6:] == tokens.Refresh[len(tokens.Refresh)-6:]
}
