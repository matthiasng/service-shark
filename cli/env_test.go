package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExpandArguments(t *testing.T) {
	require := require.New(t)

	os.Setenv("_ENV_VAL_", "123")
	result := ExpandArguments([]string{
		"normal",
		"-test", "env:_ENV_VAL_",
	})

	require.Equal(result, []string{
		"normal",
		"-test", "123",
	})
}
