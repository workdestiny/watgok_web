package app

import (
	"context"

	"github.com/moonrhythm/session"
)

var (
	// SessionName is sess
	SessionName = "sess"
)

func getSession(ctx context.Context) *session.Session {
	s, err := session.Get(ctx, SessionName)
	must(err)
	return s
}
