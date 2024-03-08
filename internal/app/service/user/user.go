package user

import (
	"context"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/internal/app/repo/postgresql/user"
	"go.uber.org/dig"
)

type (
	UserServiceImpl struct {
		dig.In
		UserRepo user.UserRepository
	}

	UserService interface {
		CreateUser(ctx context.Context, request dto.CreateUserRequest) (user entity.User, err error)
	}
)
