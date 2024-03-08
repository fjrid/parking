package user

import (
	"context"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/fjrid/parking/internal/app/model/constant"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"golang.org/x/crypto/bcrypt"
)

// @ctor
func NewUserSvc(impl UserServiceImpl) UserService {
	return &impl
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, request dto.CreateUserRequest) (user entity.User, err error) {
	isValid, err := govalidator.ValidateStruct(request)
	if err != nil {
		return
	}

	if !isValid {
		err = errors.New("data request invalid")
		return
	}

	users, err := s.UserRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return
	}

	if len(users) > 0 {
		err = errors.New("email already registered")
		return
	}

	userInterface := ctx.Value(constant.CONTEXT_USER_KEY)
	if userInterface != nil {
		user := userInterface.(dto.JWTClaims)

		request.CreatedBy = user.Email
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 1)
	if err != nil {
		return
	}

	request.Password = string(password)

	user, err = s.UserRepo.CreateUser(ctx, request)
	if err != nil {
		return
	}

	return
}
