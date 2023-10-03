package notification

// HookResult is the result of a hook execution
type HookResult interface {
	ID() uint64
	Name() string
	Logs(filter string) string
	StatusLabel() string
	Err() error
}
