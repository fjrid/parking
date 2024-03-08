package restkit_test

import (
	"errors"
	"testing"

	"github.com/fjrid/parking/pkg/restkit"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	testcases := []struct {
		TestName string
		restkit.HealthMap
		Expected   map[string]string
		ExpectedOk bool
	}{
		{
			HealthMap: restkit.HealthMap{
				"postgres": nil,
				"redis":    nil,
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"postgres": "OK",
				"redis":    "OK",
			},
		},
		{
			HealthMap: restkit.HealthMap{
				"postgres": errors.New("postgres-error"),
				"redis":    errors.New("redis-error"),
			},
			ExpectedOk: false,
			Expected: map[string]string{
				"postgres": "postgres-error",
				"redis":    "redis-error",
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			status, ok := tt.Status()
			require.Equal(t, tt.Expected, status)
			require.Equal(t, tt.ExpectedOk, ok)
		})
	}
}
