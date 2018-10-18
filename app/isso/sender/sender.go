package sender

// Sender can notify admin.
type Sender interface {
	Notify(notice string)
}