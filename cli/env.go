package cli

import (
	"os"
	"strings"
)

// ExpandArguments replace all "env:" arguments
func ExpandArguments(args []string) []string {
	result := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "env:") {
			result = append(result, os.Getenv(strings.TrimLeft(arg, "env:")))
		} else {
			result = append(result, arg)
		}
	}

	return result
}
