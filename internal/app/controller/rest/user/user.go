package user

import (
	"net/http"

	"github.com/fjrid/parking/internal/app/infra"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/service/user"
	"github.com/fjrid/parking/pkg/echokit"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	// BookCntrl is controller to book repo
	UserRestCtrl struct {
		dig.In
		UserSvc user.UserService
	}
)

var _ echokit.Router = (*UserRestCtrl)(nil)

// SetRoute to define API Route
func (c *UserRestCtrl) SetRoute(e echokit.Server) {
	r := e.Group("/users")
	r.Use(infra.JWTMiddleware())

	r.POST("", c.Create, infra.AdminMiddleware)
}

func (c *UserRestCtrl) Create(ec echo.Context) (err error) {
	var request dto.CreateUserRequest
	if err = ec.Bind(&request); err != nil {
		return err
	}

	ctx := ec.Request().Context()

	_, err = c.UserSvc.CreateUser(ctx, request)
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
