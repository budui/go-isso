package isso

import (
	"fmt"
	"runtime"
	"strings"

	"wrong.wang/x/go-isso/version"
)

// HandlerError only used in api handler
// HandlerError record origin error,
// but only show a short description to user
type HandlerError struct {
	caller      string
	description string
	origin      error
}

func (he HandlerError) Error() string {
	return he.description
}

// Unwrap return the reason cause handler broken
func (he HandlerError) Unwrap() error {
	if he.origin == nil {
		return nil
	}
	return fmt.Errorf("%s - %s - %w", he.caller, he.description, he.origin)
}

func newHandlerError(origin error, desc string) HandlerError {
	var caller string
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		caller = "unkown"
	} else {
		fn := runtime.FuncForPC(pc)
		caller = strings.TrimPrefix(fn.Name(), version.Mod)
	}
	if desc == "" {
		desc = "failed"
	}
	return HandlerError{
		description: desc,
		origin:      origin,
		caller:      caller,
	}
}
