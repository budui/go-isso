package database

import (
	"context"
	"fmt"
	"time"

	"wrong.wang/x/go-isso/isso"
)

// NewCommentGuard limit comment
func (d *Database) NewCommentGuard(ctx context.Context, c isso.Comment, uri string,
	ratelimit int, directreply int, replytoself bool, maxage int) (bool, string) {
	var n int
	d.DB.QueryRowContext(ctx, d.statement["comment_guard_ratelimit"],
		c.RemoteAddr, float64(time.Now().UnixNano())/float64(1e9)).Scan(&n)
	if n > ratelimit {
		return false, fmt.Sprintf("%s ratelimit exceeded: %d comments in 60s", c.RemoteAddr, n)
	}

	if c.Parent == nil {
		d.DB.QueryRowContext(ctx, d.statement["comment_guard_3_direct_comment"], uri, c.RemoteAddr).Scan(&n)
		if n > directreply {
			return false, fmt.Sprintf("%d direct responses to %s", n, uri)
		}
	} else if !replytoself {
		d.DB.QueryRowContext(ctx, d.statement["comment_guard_reply_to_self"],
			c.RemoteAddr, *c.Parent, float64(time.Now().UnixNano())/float64(1e9), maxage).Scan(&n)
		if n > 0 {
			return false, "edit time frame is still open"
		}
	}
	return true, ""
}
