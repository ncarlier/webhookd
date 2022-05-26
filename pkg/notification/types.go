package notification

type HookResult interface {
	ID() uint64
	Name() string
	Logs(filter string) string
	StatusLabel() string
	Err() error
}
