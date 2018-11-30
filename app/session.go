package app

import (
	"context"

	"github.com/moonrhythm/session"
)

var (
	// SessionName is sess
	SessionName = "sess_watgok"
)

func getSession(ctx context.Context) *session.Session {
	s, err := session.Get(ctx, SessionName)
	must(err)
	return s
}

func setSession(ctx context.Context, userID string) {
	s := getSession(ctx)
	s.Set("userid", userID)
}

func removeSession(ctx context.Context) {
	s := getSession(ctx)
	s.Del("userid")
}

func getUserID(ctx context.Context) string {
	return getSession(ctx).GetString("userid")
}
