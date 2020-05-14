package database

import (
	"context"

	"wrong.wang/x/go-isso/isso"
)

// Guard limit comment
func (d *Database) Guard(ctx context.Context, c isso.Comment, ratelimit int, directreply int, replytoself bool) (bool, string) {
	
	return true, ""
}
