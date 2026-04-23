package services

import (
	"context"
	"errors"
	"time"

	"github.com/RafayKhattak/aegis-iam-backend/internal/models"
	"github.com/RafayKhattak/aegis-iam-backend/internal/repository"
	"github.com/RafayKhattak/aegis-iam-backend/internal/repository/db"
	"github.com/RafayKhattak/aegis-iam-backend/pkg/hash"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type UserService interface {
	Register(ctx context.Context, req models.RegisterRequest) (models.UserResponse, error)
}

type userService struct {
	store *repository.Store
}

func NewUserService(store *repository.Store) UserService {
	return &userService{store: store}
}

func (s *userService) Register(ctx context.Context, req models.RegisterRequest) (models.UserResponse, error) {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
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
