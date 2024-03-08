package echokit_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fjrid/parking/pkg/echokit"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	var out strings.Builder
	echokit.SetRoute(nil,
		echokit.NewRouter(func(echokit.Server) { fmt.Fprintln(&out, "1") }),
		echokit.NewRouter(func(echokit.Server) { fmt.Fprintln(&out, "2") }),
	)
	require.Equal(t, "1\n2\n", out.String())
}
