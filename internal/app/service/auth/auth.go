package auth

import (
	"context"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/internal/app/repo/postgresql/user"
	"go.uber.org/dig"
)

type (
	AuthServiceImpl struct {
		dig.In
		UserRepo user.UserRepository
	}

	AuthService interface {
		Login(ctx context.Context, request dto.LoginRequest) (result dto.LoginResponse, err error)
		RegisterUser(ctx context.Context, request dto.RegisterUserRequest) (user entity.User, err error)
	}
)
