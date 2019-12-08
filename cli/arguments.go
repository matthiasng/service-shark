package cli

// Arguments container
type Arguments struct {
	Name             string
	WorkingDirectory string
	LogDirectory     string
	Command          string
	CommandArguments []string
}
