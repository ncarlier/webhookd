package worker

// WorkRequest is a request of work for a worker
type WorkRequest struct {
	Name        string
	Script      string
	Payload     string
	Args        []string
	MessageChan chan []byte
}
