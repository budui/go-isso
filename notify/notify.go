package notify

import "wrong.wang/x/go-isso/event"

// Notifier register handlers to *event.Bus
type Notifier interface {
	Register(*event.Bus)
}
