package sqkit_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/fjrid/parking/pkg/sqkit"
	"github.com/stretchr/testify/require"
)

func TestNewDeleteOption(t *testing.T) {
	expected := sq.Delete("")
	deleteOpt := sqkit.NewDeleteOption(func(sq.DeleteBuilder) sq.DeleteBuilder {
		return expected
	})
	require.Equal(t, expected, deleteOpt.CompileDelete(sq.Delete("")))
}
