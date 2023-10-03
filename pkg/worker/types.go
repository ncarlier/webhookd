package worker

// ChanWriter is a simple writer to a channel of byte.
type ChanWriter struct {
	ByteChan chan []byte
}

func (c *ChanWriter) Write(p []byte) (int, error) {
	c.ByteChan <- p
	return len(p), nil
}

// Work is a dispatched work given to a worker
type Work interface {
	ID() uint64
	Name() string
	Run() error
	Close()
	SendMessage(message string)
	Logs(filter string) string
	StatusLabel() string
	Err() error
}
