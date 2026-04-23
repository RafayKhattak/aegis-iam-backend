package services

import (
	"context"
	"testing"
	"time"

	"github.com/RafayKhattak/aegis-iam-backend/internal/models"
	"github.com/RafayKhattak/aegis-iam-backend/internal/repository/db"
	appJWT "github.com/RafayKhattak/aegis-iam-backend/pkg/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	called := m.Called(ctx, arg)
	user, _ := called.Get(0).(db.User)
	return user, called.Error(1)
}

func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	called := m.Called(ctx, email)
	user, _ := called.Get(0).(db.User)
	return user, called.Error(1)
}

func (m *MockQuerier) GetUserByID(ctx context.Context, id pgtype.UUID) (db.User, error) {
	called := m.Called(ctx, id)
	user, _ := called.Get(0).(db.User)
	return user, called.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockQuerier := &MockQuerier{}
	tokenManager := appJWT.NewTokenManager("test-secret")
	service := NewUserService(mockQuerier, tokenManager, time.Hour)

	now := time.Now().UTC().Round(time.Microsecond)
	userID := uuid.New()
	var pgUUID pgtype.UUID
	copy(pgUUID.Bytes[:], userID[:])
	pgUUID.Valid = true

	fakeUser := db.User{
		ID:           pgUUID,
		Email:        "qa@aegis.com",
		PasswordHash: "hashed",
		Role:         "user",
		CreatedAt: pgtype.Timestamptz{
			Time:  now,
			Valid: true,
		},
	}

	mockQuerier.
		On("CreateUser", mock.Anything, mock.MatchedBy(func(arg db.CreateUserParams) bool {
			if arg.Email != "qa@aegis.com" || arg.Role != "user" {
				return false
			}
			return arg.PasswordHash != "" && arg.PasswordHash != "enterprise_secure_123"
		})).
		Return(fakeUser, nil).
		Once()

	resp, err := service.Register(context.Background(), models.RegisterRequest{
		Email:    "qa@aegis.com",
		Password: "enterprise_secure_123",
	})

	require.NoError(t, err)
	require.Equal(t, userID.String(), resp.ID)
	require.Equal(t, fakeUser.Email, resp.Email)
	require.Equal(t, fakeUser.Role, resp.Role)
	require.WithinDuration(t, now, resp.CreatedAt, time.Second)
	mockQuerier.AssertExpectations(t)
}

func TestRegister_EmailExists(t *testing.T) {
	mockQuerier := &MockQuerier{}
	tokenManager := appJWT.NewTokenManager("test-secret")
	service := NewUserService(mockQuerier, tokenManager, time.Hour)

	mockQuerier.
		On("CreateUser", mock.Anything, mock.Anything).
		Return(db.User{}, &pgconn.PgError{Code: "23505"}).
		Once()

	_, err := service.Register(context.Background(), models.RegisterRequest{
		Email:    "qa@aegis.com",
		Password: "enterprise_secure_123",
	})

	require.ErrorIs(t, err, ErrEmailAlreadyExists)
	mockQuerier.AssertExpectations(t)
}
