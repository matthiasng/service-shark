package service

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FixWorkingDirectory(t *testing.T) {
	require := require.New(t)

	expectedDir := filepath.Dir(os.Args[0])

	usr, err := user.Current()
	require.NoError(err)
	require.NotEqual(expectedDir, usr.HomeDir)

	err = os.Chdir(usr.HomeDir)
	require.NoError(err)

	err = FixWorkingDirectory()
	require.NoError(err)

	newDir, err := os.Getwd()
	require.NoError(err)
	require.Equal(expectedDir, newDir)
}
