package infra

import (
	"time"

	"github.com/fjrid/parking/pkg/logruskit"
	"github.com/labstack/echo/v4"
)

type (
	// EchoCfg application configuration
	// @envconfig (prefix:"APP")
	EchoCfg struct {
		Address      string        `envconfig:"ADDRESS" default:":8089" required:"true"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
		Debug        bool          `envconfig:"DEBUG" default:"true"`
	}
)

// NewEcho return new instance of server
// @ctor
func NewEcho(cfg *EchoCfg) *echo.Echo {
	e := echo.New()
	logger := SetLogger(cfg.Debug)

	e.HideBanner = true
	e.Debug = cfg.Debug
	e.Logger = logruskit.EchoLogger(logger)
	return e
}
