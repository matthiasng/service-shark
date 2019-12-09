package service

import (
	"os"
	"path/filepath"
)

// FixWorkingDirectory changes the working directory to the exeutable directory.
// The working directory for a Windows Service is C:\Windows\System32 ...
func FixWorkingDirectory() error {
	dir := filepath.Dir(os.Args[0])
	return os.Chdir(dir)
}
