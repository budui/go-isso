package notify

import (
	"context"
	"fmt"
	"github.com/RayHY/go-isso/internal/pkg/db"
)

// Destination defines interface for a given destination service, like telegram, email and so on
type Destination interface {
	fmt.Stringer
	Send(ctx context.Context, req request) error
}

type request struct {
	comment db.Comment
	parent  db.Comment
}

// Service delivers notifications to multiple destinations
type Service struct {
	destinations []Destination
	queue        chan request

	closed uint32 // non-zero means closed. uses uint instead of bool for atomic
	ctx    context.Context
	cancel context.CancelFunc
}
