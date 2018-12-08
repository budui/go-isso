package notify

// Sender can notify admin.
type Sender interface {
	Notify(notice string)
}
