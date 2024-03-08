package infra

import (
	"context"
	"net/http"

	"github.com/fjrid/parking/internal/app/model/constant"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(dto.JWTClaims)
		},
		SigningKey: []byte(constant.JWT_CLAIM_TOKEN),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: dto.MessageUnauthorized.Translate(),
				Error: &dto.Error{
					Message: "Unauthorized",
				},
			})
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user")
			if user != nil {
				dataUser := user.(*jwt.Token)
				claims := dataUser.Claims.(*dto.JWTClaims)

				if claims != nil {
					ctx := c.Request().Context()
					newCtx := context.WithValue(ctx, constant.CONTEXT_USER_KEY, *claims)

					req := c.Request().WithContext(newCtx)
					c.SetRequest(req)
				}
			}
		},
	}

	return echojwt.WithConfig(config)
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user dto.JWTClaims
		userInterface := c.Request().Context().Value(constant.CONTEXT_USER_KEY)
		if userInterface != nil {
			user = userInterface.(dto.JWTClaims)
		}

		if user.Role != constant.ROLE_ADMIN {
			return c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: dto.MessageUnauthorized.Translate(),
				Error: &dto.Error{
					Message: "Unauthorized",
				},
			})
		}

		// current handler
		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}
