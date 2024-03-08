package auth

import (
	"context"
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fjrid/parking/internal/app/model/constant"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// @ctor
func NewAuthSvc(impl AuthServiceImpl) AuthService {
	return &impl
}

func (s *AuthServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (result dto.LoginResponse, err error) {
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

	if len(users) == 0 {
		err = errors.New("user is not found")
		return
	}

	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		err = errors.New("password is invalid")
		return
	}

	claim := dto.JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(constant.JWT_CLAIM_TOKEN))
	if err != nil {
		return
	}

	result.Email = user.Email
	result.Token = token

	return
}

func (s *AuthServiceImpl) RegisterUser(ctx context.Context, request dto.RegisterUserRequest) (user entity.User, err error) {
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

	user, err = s.UserRepo.CreateUser(ctx, dto.CreateUserRequest{
		Email:     request.Email,
		Password:  request.Password,
		Role:      constant.ROLE_USER,
		CreatedBy: request.Email,
	})
	if err != nil {
		return
	}

	return
}
