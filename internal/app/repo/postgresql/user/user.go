package user

import (
	"context"
	"database/sql"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"go.uber.org/dig"
)

type (
	UserRepositoryImpl struct {
		dig.In
		*sql.DB `name:"pg"`
	}

	// @mock
	UserRepository interface {
		CreateUser(ctx context.Context, data dto.CreateUserRequest) (res entity.User, err error)
		FindByEmail(ctx context.Context, email string) (results []entity.User, err error)
	}
)
