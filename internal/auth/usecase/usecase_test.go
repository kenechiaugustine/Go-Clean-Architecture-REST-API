package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AleksK1NG/api-mc/config"
	"github.com/AleksK1NG/api-mc/internal/auth/mock"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/AleksK1NG/api-mc/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthUC_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, nil, nil)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	ctx := context.Background()

	mockAuthRepo.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockAuthRepo.EXPECT().Register(ctx, gomock.Eq(user)).Return(user, nil)

	createdUSer, err := authUC.Register(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUSer)
	require.Nil(t, err)
}

func TestAuthUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()

	mockAuthRepo.EXPECT().Update(ctx, gomock.Eq(user)).Return(user, nil)
	mockRedisRepo.EXPECT().DeleteUserCtx(ctx, key).Return(nil)

	updatedUser, err := authUC.Update(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Nil(t, err)
}

func TestAuthUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()

	mockAuthRepo.EXPECT().Delete(ctx, gomock.Eq(user.UserID)).Return(nil)
	mockRedisRepo.EXPECT().DeleteUserCtx(ctx, key).Return(nil)

	err := authUC.Delete(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestAuthUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()

	mockRedisRepo.EXPECT().GetByIDCtx(ctx, key).Return(nil, nil)
	mockAuthRepo.EXPECT().GetByID(ctx, gomock.Eq(user.UserID)).Return(user, nil)
	mockRedisRepo.EXPECT().SetUserCtx(ctx, key, cacheDuration, user).Return(nil)

	u, err := authUC.GetByID(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}

func TestAuthUC_FindByName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil)

	userName := "name"
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}
	ctx := context.Background()

	usersList := &models.UsersList{}

	mockAuthRepo.EXPECT().FindByName(ctx, gomock.Eq(userName), query).Return(usersList, nil)

	userList, err := authUC.FindByName(ctx, userName, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userList)
}

func TestAuthUC_GetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil)

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}
	ctx := context.Background()
	usersList := &models.UsersList{}

	mockAuthRepo.EXPECT().GetUsers(ctx, query).Return(usersList, nil)

	users, err := authUC.GetUsers(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, users)
}

func TestAuthUC_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, nil)

	ctx := context.Background()
	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Email:    "email@gmail.com",
		Password: string(hashPassword),
	}

	mockAuthRepo.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := authUC.Login(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}

func TestAuthUC_UploadAvatar(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
	}

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	mockAWSRepo := mock.NewMockAWSRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, mockAWSRepo)

	ctx := context.Background()
	file := models.UploadInput{}
	uploadInfo := &minio.UploadInfo{}
	userUID := uuid.New()

	user := &models.User{
		UserID:   userUID,
		Password: "123456",
		Email:    "email@gmail.com",
	}

	mockAWSRepo.EXPECT().PutObject(ctx, gomock.Eq(file)).Return(uploadInfo, nil)
	mockAuthRepo.EXPECT().Update(ctx, gomock.Any()).Return(user, nil)

	updatedUser, err := authUC.UploadAvatar(ctx, userUID, file)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, updatedUser)
}
