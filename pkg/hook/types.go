package hook

// Status is the status of a hook
type Status int

const (
	// Idle means that the hook is not yet started
	Idle Status = iota
	// Running means that the hook is running
	Running
	// Success means that the hook over
	Success
	// Error means that the hook is over but in error
	Error
)

// Request is a hook request
type Request struct {
	Name      string
	Method    string
	Payload   string
	Args      []string
	Timeout   int
	BaseDir   string
	OutputDir string
}
