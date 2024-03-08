package dto

import "github.com/golang-jwt/jwt/v5"

type (
	JWTClaims struct {
		ID    int64  `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`

		jwt.RegisteredClaims
	}
)

type (
	LoginRequest struct {
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"required"`
	}

	LoginResponse struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	RegisterUserRequest struct {
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"required"`

		CreatedBy string `json:"-"`
	}
)
