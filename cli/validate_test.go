package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Validate(t *testing.T) {
	require := require.New(t)

	args := Arguments{}
	require.Error(Validate(args))

	args.Name = "123"
	require.Error(Validate(args))

	args.WorkingDirectory = "C:/"
	require.Error(Validate(args))

	args.Command = "_unknown_"
	require.NoError(Validate(args))
}
