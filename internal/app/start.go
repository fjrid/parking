package app

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/fjrid/parking/internal/app/controller/rest/auth"
	"github.com/fjrid/parking/internal/app/controller/rest/parking"
	"github.com/fjrid/parking/internal/app/controller/rest/user"
	"github.com/fjrid/parking/internal/app/infra"
	"github.com/fjrid/parking/pkg/echokit"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/dig"

	// enable `/debug/vars`
	_ "expvar"

	// enable `/debug/pprof` API
	_ "net/http/pprof"

	_ "github.com/fjrid/parking/docs"
)

const (
	healthCheckPath = "/application/health"
)

// Start app
func Start(
	di *dig.Container,
	cfg *infra.EchoCfg,
	e *echo.Echo,
) (err error) {
	if err := di.Invoke(SetMiddleware); err != nil {
		return err
	}
	if err := di.Invoke(SetRoute); err != nil {
		return err
	}
	if cfg.Debug {
		routes := echokit.DumpEcho(e)
		logrus.Debugf("Print routes:\n  %s\n\n", strings.Join(routes, "\n  "))
	}
	return e.StartServer(&http.Server{
		Addr:         cfg.Address,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}

// SetMiddleware set middleware
func SetMiddleware(e *echo.Echo) {
	e.Use(infra.LogMiddleware)
	e.Use(middleware.Recover())
}

// SetRoute set route
func SetRoute(
	e *echo.Echo,
	hc HealthCheck,
	userCtrl user.UserRestCtrl,
	authCtrl auth.AuthRestCtrl,
	parkingCtrl parking.ParkingRestCtrl,
) {

	// set route
	echokit.SetRoute(e, &userCtrl)
	echokit.SetRoute(e, &authCtrl)
	echokit.SetRoute(e, &parkingCtrl)

	// profiling
	e.GET(healthCheckPath, hc.Handle)
	e.HEAD(healthCheckPath, hc.Handle)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/debug/*/*", echo.WrapHandler(http.DefaultServeMux))
}
