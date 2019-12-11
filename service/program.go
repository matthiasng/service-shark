package service

// Program interface contains Start and Stop methods which are called
// when the service is started and stopped.
//
// The Start methods must be non-blocking.
//
// Implement this interface and pass it to the Run function to start your program.
type Program interface {
	// Start must be non-blocking.
	Start(Environment) error

	// Stop is called in response to syscall.SIGINT, syscall.SIGTERM, or when a
	// Windows Service is stopped.
	Stop() error

	// Name returns the service name.
	Name() string
}
