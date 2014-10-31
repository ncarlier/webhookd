package worker

type WorkRequest struct {
	Name   string
	Action string
	Args   []string
}
