package auth

import (
	"net/http"

	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/service/auth"
	"github.com/fjrid/parking/pkg/echokit"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	// BookCntrl is controller to book repo
	AuthRestCtrl struct {
		dig.In
		AuthSvc auth.AuthService
	}
)

var _ echokit.Router = (*AuthRestCtrl)(nil)

// SetRoute to define API Route
func (c *AuthRestCtrl) SetRoute(e echokit.Server) {
	r := e.Group("/auth")
	r.POST("/login", c.Login)
	r.POST("/register", c.Register)
}

func (c *AuthRestCtrl) Login(ec echo.Context) (err error) {
	var request dto.LoginRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	ctx := ec.Request().Context()

	result, err := c.AuthSvc.Login(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
		Data:    result,
	})
}

func (c *AuthRestCtrl) Register(ec echo.Context) (err error) {
	var request dto.RegisterUserRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	ctx := ec.Request().Context()

	_, err = c.AuthSvc.RegisterUser(ctx, request)
	if err != nil {
		return ec.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: dto.MessageBadRequest.Translate(),
			Error: &dto.Error{
				Message: err.Error(),
			},
		})
	}

	return ec.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: dto.MessageSucessfully.Translate(),
	})
}
