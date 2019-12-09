package service

// Environment contains information about the environment
// your application is running in.
type Environment interface {
	// IsWindowsService reports whether the program is running as a Windows Service.
	IsWindowsService() bool

	// ExitService can be used to signal a service crash. Service will exit with a user define error code 3.
	ExitService(error)
}
