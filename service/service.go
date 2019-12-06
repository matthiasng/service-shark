package service

// Service interface contains Start and Stop methods which are called
// when the service is started and stopped.
//
// The Start methods must be non-blocking.
//
// Implement this interface and pass it to the Run function to start your program.
type Service interface {
	// Start is called after Init. This method must be non-blocking.
	Start(Environment) error

	// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
	// Windows Service is stopped.
	Stop() error
}

// Environment contains information about the environment
// your application is running in.
type Environment interface {
	// IsWindowsService reports whether the program is running as a Windows Service.
	IsWindowsService() bool

	// ExitService can be used to signal a service crash. Service will exit with a user define error code 3.
	ExitService(error)

	// Name returns the service name.
	Name() string
}
