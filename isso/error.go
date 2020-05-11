package isso

// HandlerError only used in api handler
// HandlerError record origin error, 
// but only show a short description to user
type HandlerError struct {
	handlerName string
	description string
	origin error
}

func (he HandlerError) Error() string {
	return he.description
}

// Unwrap return the reason cause handler broken
func (he HandlerError) Unwrap() error {
	return he.origin
}

func newHandlerError(desc string, origin error) HandlerError {
	return HandlerError{
		description: desc,
		origin: origin,
	}
}