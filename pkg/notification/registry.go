package notification

import (
	"fmt"
	"net/url"
)

// NotifierCreator function for create a notifier
type NotifierCreator func(uri *url.URL) (Notifier, error)

// Registry of all Notifiers
var registry = map[string]NotifierCreator{}

// Register a Notifier to the registry
func Register(scheme string, creator NotifierCreator) {
	registry[scheme] = creator
}

// NewNotifier create new Notifier
func NewNotifier(uri string) (Notifier, error) {
	if uri == "" {
		return nil, nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid notification URL: %s", uri)
	}
	creator, ok := registry[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("unsupported notification scheme: %s", u.Scheme)
	}
	return creator(u)
}
