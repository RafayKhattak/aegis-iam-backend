package services

import (
	"context"
	"errors"
	"time"

	"github.com/RafayKhattak/aegis-iam-backend/internal/models"
	"github.com/RafayKhattak/aegis-iam-backend/internal/repository/db"
	"github.com/RafayKhattak/aegis-iam-backend/pkg/hash"
	appJWT "github.com/RafayKhattak/aegis-iam-backend/pkg/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService interface {
	Register(ctx context.Context, req models.RegisterRequest) (models.UserResponse, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	querier       db.Querier
	tokenManager  *appJWT.TokenManager
	tokenDuration time.Duration
}

func NewUserService(querier db.Querier, tokenManager *appJWT.TokenManager, tokenDuration time.Duration) UserService {
	return &userService{
		querier:       querier,
		tokenManager:  tokenManager,
		tokenDuration: tokenDuration,
	}
}

func (s *userService) Register(ctx context.Context, req models.RegisterRequest) (models.UserResponse, error) {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.querier.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return models.UserResponse{}, ErrEmailAlreadyExists
		}
		return models.UserResponse{}, err
	}

	createdAt := time.Time{}
	if user.CreatedAt.Valid {
		createdAt = user.CreatedAt.Time
	}

	return models.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: createdAt,
	}, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.querier.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := hash.CheckPassword(password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	if !user.ID.Valid {
		return "", errors.New("user id is invalid")
	}

	userID, err := uuid.FromBytes(user.ID.Bytes[:])
	if err != nil {
		return "", err
	}

	token, err := s.tokenManager.GenerateToken(userID, user.Email, s.tokenDuration)
	if err != nil {
		return "", err
	}

	return token, nil
}
