package cli

import (
	"errors"
)

// Validate command line arguments
func Validate(arguments Arguments) error {
	if len(arguments.Name) == 0 {
		return errors.New("Name not set (-name)")
	}

	if len(arguments.WorkingDirectory) == 0 {
		return errors.New("Working directory not set (-workdir)")
	}

	if len(arguments.Command) == 0 {
		return errors.New("Command not set (-cmd)")
	}

	return nil
}
