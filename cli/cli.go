package cli

import (
	"os"
	"strings"
)

// BindArguments replace all "bind:" arguments
func BindArguments(args []string) []string {
	result := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "bind:") {
			result = append(result, os.Getenv(strings.TrimLeft(arg, "bind:")))
		} else {
			result = append(result, arg)
		}
	}

	return result
}
